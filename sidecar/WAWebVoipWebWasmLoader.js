/**
 * WAWebVoipWebWasmLoader.js - WhatsApp Web VoIP WASM Loader
 * 
 * Bu dosya WhatsApp Web'in kendi loader script'idir.
 * Gerçek dosyayı `npm run fetch-wasm` ile indirip buraya koyun.
 * 
 * Emscripten Module factory function export eder:
 *   export default function createVoipModule(Module) { ... return Module; }
 * 
 * Module objesi şunları bekler:
 *   - locateFile(file): .wasm/.data/.js dosyalarının yolu
 *   - print(text): console.log
 *   - printErr(text): console.error
 *   - noInitialRun: true
 *   - noExitRuntime: true
 *   - onSignalingXmpp, sendDataToRelay, captureAudio, playbackAudio callbacks
 */

// PLACEHOLDER - Gerçek dosya fetch-wasm ile indirilecek
// Bu dosya WhatsApp Web'in minified JS'idir, manuel yazılamaz.

export default function createVoipModule(Module) {
  console.error('[WAWebVoipWebWasmLoader] PLACEHOLDER - Real loader not loaded!');
  console.error('[WAWebVoipWebWasmLoader] Run: npm run fetch-wasm');
  console.error('[WAWebVoipWebWasmLoader] Or manually download from WhatsApp Web and save as WAWebVoipWebWasmLoader.js');
  
  // Minimal mock for testing
  Module.initVoipStack = () => console.log('[MOCK] initVoipStack called');
  Module.handleIncomingSignalingOffer = () => console.log('[MOCK] handleIncomingSignalingOffer called');
  Module.handleIncomingSignalingMessage = () => console.log('[MOCK] handleIncomingSignalingMessage called');
  Module.acceptCall = () => console.log('[MOCK] acceptCall called');
  Module.rejectCall = () => console.log('[MOCK] rejectCall called');
  
  return Module;
}