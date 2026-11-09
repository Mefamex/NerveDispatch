/**
 * WhatsApp Web VoIP WASM Sidecar
 * 
 * Bu sidecar:
 * 1. Go bridge (whatsmeow) ile stdin/stdout üzerinden JSON line IPC kurar
 * 2. WhatsApp Web VoIP WASM'i (Emscripten module) yükler
 * 3. Gelen/iden call sinyalleri, relay paketleri, audio I/O callback'lerini yönetir
 * 4. Go bridge'e ham call stanza'ları iletir
 * 
 * IPC Protocol (JSON Lines):
 * Go -> Sidecar:
 *   { "type": "call_stanza", "stanza": {...} }           // Ham call node (offer/transport/relaylatency/terminate)
 *   { "type": "accept_call", "callId": "...", "audio": true, "video": false }
 *   { "type": "reject_call", "callId": "..." }
 *   { "type": "init", "selfJid": "...", "phoneJid": "...", "lidJid": "..." }
 * 
 * Sidecar -> Go:
 *   { "type": "call_stanza_out", "stanza": {...} }       // WASM'den çıkan giden sinyal (preaccept/accept)
 *   { "type": "call_state", "callId": "...", "state": "ringing|connecting|connected|ended", "from": "..." }
 *   { "type": "audio_level", "level": 0.5 }              // Mic seviyesi (opsiyonel)
 *   { "type": "error", "message": "..." }
 *   { "type": "ready" }
 */

import { createRequire } from 'module';
import { stdin, stdout } from 'process';
import { createInterface } from 'readline';
import { VoipStack } from './voip-stack.js';
import { CallBridge } from './call-bridge.js';
import { fetchLatestWasm, getCachedWasmPath } from './wasm-loader.js';

// ==================== IPC SETUP ====================

const rl = createInterface({ input: stdin, output: stdout, terminal: false });

let voipStack = null;
let callBridge = null;
let isInitialized = false;

// Go bridge'e mesaj gönder
function sendToBridge(msg) {
  stdout.write(JSON.stringify(msg) + '\n');
}

// ==================== MESSAGE HANDLERS ====================

async function handleInit(data) {
  try {
    const { selfJid, phoneJid, lidJid } = data;
    
    sendToBridge({ type: 'status', message: 'Loading WASM...' });
    
    // WASM'i indir/yükle
    const wasmPath = await fetchLatestWasm();
    sendToBridge({ type: 'status', message: `WASM loaded: ${wasmPath}` });
    
    // VoIP Stack başlat
    voipStack = new VoipStack(wasmPath);
    await voipStack.init();
    
    // Call Bridge oluştur (signaling + relay + audio)
    callBridge = new CallBridge(voipStack, sendToBridge);
    
    // VoIP stack'i initialize et (audio drivers register olur)
    voipStack.initVoipStack(selfJid, phoneJid, lidJid);
    
    isInitialized = true;
    sendToBridge({ type: 'ready' });
    sendToBridge({ type: 'status', message: 'VoIP stack initialized' });
    
  } catch (err) {
    sendToBridge({ type: 'error', message: `Init failed: ${err.message}` });
    console.error('[Sidecar] Init error:', err);
  }
}

function handleCallStanza(data) {
  if (!isInitialized || !callBridge) {
    sendToBridge({ type: 'error', message: 'Not initialized' });
    return;
  }
  callBridge.handleIncomingStanza(data.stanza);
}

function handleAcceptCall(data) {
  if (!isInitialized || !voipStack) {
    sendToBridge({ type: 'error', message: 'Not initialized' });
    return;
  }
  voipStack.acceptCall(data.audio !== false, data.video === true);
}

function handleRejectCall(data) {
  if (!isInitialized || !voipStack) {
    sendToBridge({ type: 'error', message: 'Not initialized' });
    return;
  }
  voipStack.rejectCall();
}

// ==================== MAIN LOOP ====================

console.error('[Sidecar] Starting WhatsApp VoIP Sidecar...');

rl.on('line', (line) => {
  try {
    const msg = JSON.parse(line);
    switch (msg.type) {
      case 'init':
        handleInit(msg);
        break;
      case 'call_stanza':
        handleCallStanza(msg);
        break;
      case 'accept_call':
        handleAcceptCall(msg);
        break;
      case 'reject_call':
        handleRejectCall(msg);
        break;
      default:
        console.error(`[Sidecar] Unknown message type: ${msg.type}`);
    }
  } catch (err) {
    console.error('[Sidecar] Parse error:', err.message);
  }
});

rl.on('close', () => {
  console.error('[Sidecar] IPC closed, shutting down...');
  if (voipStack) voipStack.cleanup();
  process.exit(0);
});

process.on('SIGINT', () => {
  console.error('[Sidecar] SIGINT received');
  rl.close();
});

process.on('SIGTERM', () => {
  console.error('[Sidecar] SIGTERM received');
  rl.close();
});

console.error('[Sidecar] Ready, waiting for init command from Go bridge...');