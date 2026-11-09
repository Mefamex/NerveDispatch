/**
 * CallBridge - Call Signaling, Relay & Audio Management
 * 
 * Bu sınıf:
 * 1. Gelen ham call stanza'ları (offer/transport/relaylatency/terminate) işler
 * 2. WAP binary encode/decode yapar
 * 3. Signaling callback'lerini yönetir (preaccept/accept gönderir)
 * 4. UDP relay paketlerini yönlendirir
 * 5. Virtual audio driver'ları yönetir (silence/prerecorded/mic)
 */

import { encodeBinaryNode, decodeBinaryNode, getAllBinaryNodeKeys } from '@whiskeysockets/baileys';
import { Buffer } from 'buffer';
import dgram from 'dgram';

export class CallBridge {
  constructor(voipStack, sendToBridge) {
    this.voipStack = voipStack;
    this.sendToBridge = sendToBridge;
    
    // UDP sockets for relay
    this.udp4 = null;
    this.udp6 = null;
    this.relayMap = new Map(); // callId -> { ip, port, ipv6 }
    
    // Audio state
    this.audioDriver = 'silence'; // 'silence' | 'prerecorded' | 'mic'
    this.prerecordedBuffer = null;
    this.prerecordedOffset = 0;
    
    // Pending calls
    this.pendingCalls = new Map(); // callId -> { peerJid, platform, version, ... }
    
    this.setupVoipCallbacks();
  }

  setupVoipCallbacks() {
    this.voipStack.setCallbacks({
      onSignalingXmpp: (peerJid, callId, xmlPayload) => this.onSignalingXmpp(peerJid, callId, xmlPayload),
      sendDataToRelay: (data, ip, port) => this.onSendDataToRelay(data, ip, port),
      captureAudio: (buffer, sampleRate) => this.onCaptureAudio(buffer, sampleRate),
      playbackAudio: (buffer, sampleRate) => this.onPlaybackAudio(buffer, sampleRate),
    });
  }

  // ==================== INCOMING STANZA HANDLING ====================

  /**
   * Go bridge'den gelen ham call stanza'sını işle
   * stanza: { tag: "call", attrs: {...}, content: [...] }
   */
  handleIncomingStanza(stanza) {
    if (!stanza || stanza.tag !== 'call') {
      console.warn('[CallBridge] Not a call stanza, ignoring');
      return;
    }

    const peerJid = stanza.attrs.from;
    const callId = stanza.attrs.id;
    
    // Call content'ini bul (offer/transport/relaylatency/terminate)
    const signalingChild = stanza.content?.find(c => 
      c.tag === 'offer' || c.tag === 'transport' || c.tag === 'relaylatency' || c.tag === 'terminate'
    );

    if (!signalingChild) {
      console.warn('[CallBridge] No signaling child in call stanza');
      return;
    }

    console.log(`[CallBridge] Received ${signalingChild.tag} from ${peerJid}, callId: ${callId}`);

    switch (signalingChild.tag) {
      case 'offer':
        this.handleOffer(peerJid, callId, signalingChild);
        break;
      case 'transport':
      case 'relaylatency':
        this.handleSignalingMessage(peerJid, signalingChild);
        break;
      case 'terminate':
        this.handleTerminate(peerJid, callId);
        break;
    }
  }

  /**
   * Incoming call offer işle
   */
  handleOffer(peerJid, callId, offerNode) {
    try {
      // 1. Offer node'unu WAP binary encode et -> base64
      const wapBytes = encodeBinaryNode(offerNode);
      const stanzaB64 = Buffer.from(wapBytes).toString('base64');
      
      // 2. Offer attributes'larını al
      const platform = offerNode.attrs.platform || 'web';
      const version = offerNode.attrs.version || '0';
      const e = offerNode.attrs['call-id'] || '0';
      const t = offerNode.attrs.t || '0';
      const offline = offerNode.attrs.offline ? 1 : 0;
      const notContact = offerNode.attrs['not-contact'] ? 1 : 0;
      
      // 3. Encrypted call key varsa decrypt et (Baileys signal repo ile)
      // NOT: Bu kısım için Baileys'in signal repository'sine erişim lazım
      // Şimdilik encrypted node'u olduğu gibi gönderiyoruz, WASM parse edemezse hata verecek
      
      // 4. Pending call kaydet
      this.pendingCalls.set(callId, { peerJid, platform, version, e, t, offline, notContact });
      
      // 5. WASM'e gönder
      this.voipStack.handleOffer(stanzaB64, platform, version, e, t, offline, notContact, peerJid);
      
      // 6. Bridge'e bildir
      this.sendToBridge({ 
        type: 'call_state', 
        callId, 
        state: 'ringing', 
        from: peerJid,
        isVideo: offerNode.content?.some(c => c.tag === 'video') || false
      });
      
    } catch (err) {
      console.error('[CallBridge] Error handling offer:', err);
      this.sendToBridge({ type: 'error', message: `Offer handling failed: ${err.message}` });
    }
  }

  /**
   * Transport/relaylatency mesajlarını işle
   */
  handleSignalingMessage(peerJid, signalingNode) {
    try {
      // Sadece signaling child'ını encode et
      const wapBytes = encodeBinaryNode(signalingNode);
      const stanzaB64 = Buffer.from(wapBytes).toString('base64');
      
      this.voipStack.handleSignalingMessage(stanzaB64, peerJid);
    } catch (err) {
      console.error('[CallBridge] Error handling signaling message:', err);
    }
  }

  /**
   * Call terminate
   */
  handleTerminate(peerJid, callId) {
    this.pendingCalls.delete(callId);
    this.relayMap.delete(callId);
    this.sendToBridge({ type: 'call_state', callId, state: 'ended', from: peerJid });
  }

  // ==================== SIGNALING CALLBACK (WASM -> WHATSAPP) ====================

  /**
   * WASM'den çıkan giden sinyal (preaccept/accept/transport)
   * xmlPayload: base64 encoded WAP binary node (child node)
   */
  async onSignalingXmpp(peerJid, callId, xmlPayload) {
    try {
      // xmlPayload'yi decode et
      const payloadBytes = Buffer.from(xmlPayload, 'base64');
      const node = decodeBinaryNode(payloadBytes);
      
      // Eğer child node ise (offer/accept/transport), call wrapper'ına sar
      const callNode = this.wrapOutgoingSignalingNode(node, peerJid, callId);
      
      // Go bridge'e gönder
      this.sendToBridge({ type: 'call_stanza_out', stanza: callNode });
      
      console.log(`[CallBridge] Sent ${node.tag} for call ${callId} to ${peerJid}`);
    } catch (err) {
      console.error('[CallBridge] Error in onSignalingXmpp:', err);
    }
  }

  wrapOutgoingSignalingNode(node, peerJid, callId) {
    if (node.tag === 'call') {
      node.attrs = node.attrs || {};
      node.attrs.to = peerJid;
      node.attrs.id = node.attrs.id || callId || this.makeStanzaId();
      return node;
    }

    return {
      tag: 'call',
      attrs: {
        to: peerJid,
        id: callId || this.makeStanzaId(),
        from: '', // Go bridge doldurur
      },
      content: [node],
    };
  }

  makeStanzaId() {
    return Math.random().toString(36).substring(2, 15) + Date.now().toString(36);
  }

  // ==================== RELAY UDP HANDLING ====================

  onSendDataToRelay(data, ip, port) {
    // UDP socket hazırla
    this.ensureUdpSockets();
    
    const isIPv6 = ip.includes(':');
    const socket = isIPv6 ? this.udp6 : this.udp4;
    
    if (socket) {
      socket.send(data, port, ip, (err) => {
        if (err) console.error('[CallBridge] Relay send error:', err);
      });
    } else {
      console.warn('[CallBridge] UDP socket not ready for relay');
    }
  }

  ensureUdpSockets() {
    if (this.udp4 && this.udp6) return;
    
    // IPv4
    this.udp4 = dgram.createSocket('udp4');
    this.udp4.on('message', (packet, rinfo) => {
      this.onRelayPacket(packet, rinfo.address, rinfo.port);
    });
    this.udp4.on('error', (err) => console.error('[CallBridge] UDP4 error:', err));
    
    // IPv6
    try {
      this.udp6 = dgram.createSocket('udp6');
      this.udp6.on('message', (packet, rinfo) => {
        this.onRelayPacket(packet, rinfo.address, rinfo.port);
      });
      this.udp6.on('error', (err) => console.error('[CallBridge] UDP6 error:', err));
    } catch (err) {
      console.warn('[CallBridge] IPv6 not available:', err.message);
    }
  }

  onRelayPacket(packet, ip, port) {
    // Gelen relay paketini WASM'e ver
    // VoipStack'in handleIncomingRelayPacket export'u varsa çağır
    if (this.voipStack.module?.handleIncomingRelayPacket) {
      // WASM memory'sine kopyala ve çağır
      // Bu kısım WASM export'larına göre uyarlanmalı
      console.log(`[CallBridge] Relay rx ${packet.length} bytes from ${ip}:${port}`);
    }
  }

  // ==================== AUDIO HANDLING ====================

  /**
   * Audio capture callback - WASM mikrofon verisi istediğinde çağrılır
   * buffer: Int16Array (WASM memory'sine yazılacak)
   * sampleRate: 16000 (WhatsApp Opus standardı)
   */
  onCaptureAudio(buffer, sampleRate) {
    const frameSize = buffer.length;
    
    switch (this.audioDriver) {
      case 'silence':
        buffer.fill(0);
        break;
        
      case 'prerecorded':
        if (this.prerecordedBuffer) {
          const remaining = this.prerecordedBuffer.length - this.prerecordedOffset;
          const toCopy = Math.min(frameSize, remaining);
          
          if (toCopy > 0) {
            buffer.set(this.prerecordedBuffer.subarray(this.prerecordedOffset, this.prerecordedOffset + toCopy));
            this.prerecordedOffset += toCopy;
          }
          
          if (this.prerecordedOffset >= this.prerecordedBuffer.length) {
            // Loop veya sessizlik
            this.prerecordedOffset = 0; // Loop
            // buffer.fill(0); // veya sessizlik
          }
        } else {
          buffer.fill(0);
        }
        break;
        
      case 'mic':
        // Gerçek mikrofon için: external process (ffmpeg/arecord) stdout'dan oku
        // Şimdilik silence
        buffer.fill(0);
        break;
    }
  }

  /**
   * Audio playback callback - WASM hoparlör verisi verdiğinde çağrılır
   */
  onPlaybackAudio(buffer, sampleRate) {
    // Gelen sesi logla, dosyaya yaz, veya speaker'a çal
    // console.log(`[CallBridge] Playback ${buffer.length} frames at ${sampleRate}Hz`);
  }

  // ==================== PUBLIC API ====================

  setAudioDriver(driver, options = {}) {
    this.audioDriver = driver;
    this.prerecordedOffset = 0;
    
    if (driver === 'prerecorded' && options.filePath) {
      this.loadPrerecordedAudio(options.filePath);
    }
  }

  async loadPrerecordedAudio(filePath) {
    // PCM dosyası oku (16kHz, mono, s16le)
    const data = await fs.readFile(filePath);
    // Int16Array'e çevir (2 byte per sample)
    this.prerecordedBuffer = new Int16Array(data.buffer, data.byteOffset, data.byteLength / 2);
    console.log(`[CallBridge] Loaded prerecorded audio: ${this.prerecordedBuffer.length} samples`);
  }

  acceptCall(callId, audio = true, video = false) {
    const call = this.pendingCalls.get(callId);
    if (call) {
      this.voipStack.acceptCall(audio, video);
    }
  }

  rejectCall(callId) {
    const call = this.pendingCalls.get(callId);
    if (call) {
      this.voipStack.rejectCall();
    }
  }

  cleanup() {
    if (this.udp4) this.udp4.close();
    if (this.udp6) this.udp6.close();
  }
}

// fs import
import { promises as fs } from 'fs';