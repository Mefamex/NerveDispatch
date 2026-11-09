/**
 * VoipStack - WhatsApp Web VoIP WASM Wrapper
 * 
 * Emscripten Module'ü yükler ve VoIP stack export'larını expose eder.
 * Blog yazısındaki (Salman Ansari) yöntemi baz alınarak yazılmıştır.
 */

import { createRequire } from 'module';
import { promises as fs } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';

const require = createRequire(import.meta.url);
const __dirname = dirname(fileURLToPath(import.meta.url));

export class VoipStack {
  constructor(wasmPath) {
    this.wasmPath = wasmPath;
    this.module = null;
    this.initialized = false;
    this.callbacks = {
      onSignalingXmpp: null,
      sendDataToRelay: null,
      captureAudio: null,
      playbackAudio: null,
    };
  }

  async init() {
    // WhatsApp Web'in WASM loader script'ini yükle
    const loaderPath = join(__dirname, 'WAWebVoipWebWasmLoader.js');
    
    let loaderFactory;
    try {
      loaderFactory = require(loaderPath);
    } catch (err) {
      throw new Error(`WASM loader not found at ${loaderPath}. Run 'npm run fetch-wasm' first.`);
    }

    // Emscripten Module config
    const Module = {
      locateFile: (file) => {
        if (file.endsWith('.wasm')) return this.wasmPath;
        if (file.endsWith('.data')) return this.wasmPath.replace('.wasm', '.data');
        if (file.endsWith('.js')) return this.wasmPath.replace('.wasm', '.js');
        return file;
      },
      print: (text) => console.log('[VoIP:WASM]', text),
      printErr: (text) => console.error('[VoIP:WASM:ERR]', text),
      noInitialRun: true,
      noExitRuntime: true,
      thisProgram: 'wa-voip',
      preRun: [
        () => {
          // Virtual FS için persistent dizin oluştur
          if (Module.FS) {
            const persistentDir = join(process.env.HOME || '/tmp', '.wa-voip-sidecar');
            try {
              Module.FS.mkdirTree(persistentDir);
            } catch {}
          }
        },
      ],
      // Callback'ler - WASM JS'ye çağrı yapacak
      onSignalingXmpp: (peerJid, callId, xmlPayload) => {
        if (this.callbacks.onSignalingXmpp) {
          this.callbacks.onSignalingXmpp(peerJid, callId, xmlPayload);
        }
      },
      sendDataToRelay: (dataPtr, len, ip, port) => {
        if (this.callbacks.sendDataToRelay) {
          // dataPtr: WASM memory'sindeki pointer, len: byte length
          const data = new Uint8Array(Module.HEAPU8.buffer, dataPtr, len);
          this.callbacks.sendDataToRelay(data, ip, port);
        }
      },
      // Audio capture callback - WASM ses verisi isteyecek
      captureAudio: (bufferPtr, numFrames, sampleRate) => {
        if (this.callbacks.captureAudio) {
          const buffer = new Int16Array(Module.HEAPU8.buffer, bufferPtr, numFrames);
          this.callbacks.captureAudio(buffer, sampleRate);
        }
      },
      // Audio playback callback - WASM ses verisi oynatacak
      playbackAudio: (bufferPtr, numFrames, sampleRate) => {
        if (this.callbacks.playbackAudio) {
          const buffer = new Int16Array(Module.HEAPU8.buffer, bufferPtr, numFrames);
          this.callbacks.playbackAudio(buffer, sampleRate);
        }
      },
    };

    // Module'i oluştur
    this.module = loaderFactory(Module);
    
    // Module ready olana kadar bekle
    await new Promise((resolve, reject) => {
      if (this.module.then) {
        this.module.then(resolve).catch(reject);
      } else {
        // Sync module
        resolve();
      }
    });

    // Export'ları bind et
    this.bindExports();
    this.initialized = true;
    console.log('[VoipStack] WASM module loaded and initialized');
  }

  bindExports() {
    const m = this.module;
    if (!m) throw new Error('Module not loaded');

    // VoIP Stack functions
    this.initVoipStack = m.initVoipStack || m._initVoipStack;
    this.handleIncomingSignalingOffer = m.handleIncomingSignalingOffer || m._handleIncomingSignalingOffer;
    this.handleIncomingSignalingMessage = m.handleIncomingSignalingMessage || m._handleIncomingSignalingMessage;
    this.acceptCall = m.acceptCall || m._acceptCall;
    this.rejectCall = m.rejectCall || m._rejectCall;
    this.terminateCall = m.terminateCall || m._terminateCall;
    this.getCallState = m.getCallState || m._getCallState;
    this.setAudioDevices = m.setAudioDevices || m._setAudioDevices;

    if (!this.initVoipStack) {
      console.warn('[VoipStack] Warning: initVoipStack not found in exports');
      console.log('[VoipStack] Available exports:', Object.keys(m).filter(k => k.includes('Voip') || k.includes('Call') || k.includes('call')));
    }
  }

  // ==================== HIGH-LEVEL API ====================

  /**
   * VoIP stack'i başlat
   * @param {string} selfJid - Kendi JID'imiz (örn: "905551234567@s.whatsapp.net")
   * @param {string} phoneJid - Telefon JID (LID olmayan)
   * @param {string} lidJid - LID JID
   * @param {boolean} isVideo - Video call mi
   */
  initVoipStackWrapper(selfJid, phoneJid, lidJid, isVideo = false) {
    if (!this.initVoipStack) throw new Error('initVoipStack not available');
    return this.initVoipStack(selfJid, phoneJid, lidJid, isVideo);
  }

  /**
   * Gelen call offer'ını işle
   * @param {string} stanzaB64 - Base64 encoded WAP binary node (offer child)
   * @param {string} platform - "web" | "android" | "ios"
   * @param {string} version - Call version
   * @param {string} e - Ephemeral key ID
   * @param {string} t - Timestamp
   * @param {number} offline - Offline mi
   * @param {number} notContact - Contact değil mi
   * @param {string} peerJid - Karşı taraf JID
   */
  handleOffer(stanzaB64, platform, version, e, t, offline, notContact, peerJid) {
    if (!this.handleIncomingSignalingOffer) throw new Error('handleIncomingSignalingOffer not available');
    return this.handleIncomingSignalingOffer(
      stanzaB64,
      platform,
      version,
      String(e || 0),
      String(t || 0),
      offline ? 1 : 0,
      notContact ? 1 : 0,
      peerJid,
      null
    );
  }

  /**
   * Gelen transport/relaylatency mesajlarını işle
   */
  handleSignalingMessage(stanzaB64, peerJid) {
    if (!this.handleIncomingSignalingMessage) throw new Error('handleIncomingSignalingMessage not available');
    return this.handleIncomingSignalingMessage(stanzaB64, peerJid, null);
  }

  /**
   * Aramayı kabul et
   * @param {boolean} audio - Sesli arama
   * @param {boolean} video - Görüntülü arama
   */
  acceptCall(audio = true, video = false) {
    if (!this.acceptCall) throw new Error('acceptCall not available');
    return this.acceptCall(audio, video);
  }

  /**
   * Aramayı reddet
   */
  rejectCall() {
    if (!this.rejectCall) throw new Error('rejectCall not available');
    return this.rejectCall();
  }

  /**
   * Callback'leri ayarla
   */
  setCallbacks(callbacks) {
    this.callbacks = { ...this.callbacks, ...callbacks };
  }

  cleanup() {
    if (this.module && this.module.FS) {
      try {
        this.module.FS.unmount('/persistent');
      } catch {}
    }
    this.initialized = false;
  }
}