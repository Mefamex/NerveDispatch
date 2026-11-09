/**
 * WASM Loader - WhatsApp Web VoIP WASM'ini indirir ve cache'ler
 * 
 * WASM URL'si her WhatsApp Web deploy'unda değişir.
 * Pattern: https://static.whatsapp.net/rsrc.php/<path>/<hash>.wasm
 * Loader script: WAWebVoipWebWasmLoader.js (daha stabil)
 */

import { promises as fs } from 'fs';
import { join } from 'path';
import { fileURLToPath } from 'url';
import https from 'https';

const __dirname = join(process.cwd());
const CACHE_DIR = join(process.env.HOME || '/tmp', '.wa-voip-sidecar', 'cache');
const LOADER_URL = 'https://static.whatsapp.net/rsrc.php/v4/yj/r/atBZ6RzieZ4.js'; // Example - needs discovery

// Bilinen WASM hash'ler (blogdan ve geçmişten)
const KNOWN_WASM_HASHES = [
  'KwJJIha0f3H', // Blogdan (Nisan 2026)
];

export async function fetchLatestWasm() {
  await fs.mkdir(CACHE_DIR, { recursive: true });
  
  // 1. Cache'de varsa kullan
  const cached = getCachedWasmPath();
  if (cached) {
    console.log('[WASM Loader] Using cached WASM:', cached);
    return cached;
  }
  
  // 2. WhatsApp Web'ten WASM URL'sini keşfet
  const wasmUrl = await discoverWasmUrl();
  if (!wasmUrl) {
    throw new Error('Could not discover WASM URL. Try manually: npm run fetch-wasm -- --url <url>');
  }
  
  // 3. İndir
  const wasmPath = join(CACHE_DIR, `wa-voip-${Date.now()}.wasm`);
  await downloadFile(wasmUrl, wasmPath);
  
  // 4. Loader script'i de indir
  const loaderUrl = await discoverLoaderUrl();
  if (loaderUrl) {
    const loaderPath = join(__dirname, 'WAWebVoipWebWasmLoader.js');
    await downloadFile(loaderUrl, loaderPath);
  }
  
  // 5. Metadata kaydet
  await fs.writeFile(
    join(CACHE_DIR, 'latest.json'),
    JSON.stringify({ wasmUrl, wasmPath, loaderUrl, timestamp: Date.now() }, null, 2)
  );
  
  return wasmPath;
}

export function getCachedWasmPath() {
  try {
    const metaPath = join(CACHE_DIR, 'latest.json');
    const meta = JSON.parse(require('fs').readFileSync(metaPath, 'utf-8'));
    if (meta.wasmPath && require('fs').existsSync(meta.wasmPath)) {
      return meta.wasmPath;
    }
  } catch {}
  return null;
}

async function discoverWasmUrl() {
  // Yöntem 1: WhatsApp Web ana sayfasını çek, JS bundle'larda wasm referansı ara
  try {
    const html = await fetchUrl('https://web.whatsapp.com/');
    
    // JS bundle URL'lerini bul
    const jsUrls = html.match(/static\.whatsapp\.net\/rsrc\.php\/v4[a-zA-Z0-9/_-]+\.js/g);
    if (jsUrls) {
      for (const jsUrl of [...new Set(jsUrls)].slice(0, 10)) {
        const fullUrl = jsUrl.startsWith('http') ? jsUrl : `https://${jsUrl}`;
        const jsContent = await fetchUrl(fullUrl);
        
        // WASM referansı ara
        const wasmMatch = jsContent.match(/https?:\/\/static\.whatsapp\.net\/rsrc\.php\/[^"'\s]+\.wasm/);
        if (wasmMatch) {
          console.log('[WASM Loader] Discovered WASM URL:', wasmMatch[0]);
          return wasmMatch[0];
        }
        
        // Loader script referansı ara
        const loaderMatch = jsContent.match(/WAWebVoipWebWasmLoader[^"'\s]*/);
        if (loaderMatch) {
          console.log('[WASM Loader] Found loader reference:', loaderMatch[0]);
        }
      }
    }
  } catch (err) {
    console.warn('[WASM Loader] Auto-discovery failed:', err.message);
  }
  
  // Yöntem 2: Bilinen hash'leri dene
  for (const hash of KNOWN_WASM_HASHES) {
    const testUrl = `https://static.whatsapp.net/rsrc.php/ys/r/${hash}.wasm`;
    try {
      const res = await fetchHead(testUrl);
      if (res.ok) {
        console.log('[WASM Loader] Found working known hash:', testUrl);
        return testUrl;
      }
    } catch {}
  }
  
  return null;
}

async function discoverLoaderUrl() {
  // Loader script URL'sini keşfet
  try {
    const html = await fetchUrl('https://web.whatsapp.com/');
    const jsUrls = html.match(/static\.whatsapp\.net\/rsrc\.php\/v4[a-zA-Z0-9/_-]+\.js/g);
    
    if (jsUrls) {
      for (const jsUrl of [...new Set(jsUrls)].slice(0, 10)) {
        const fullUrl = jsUrl.startsWith('http') ? jsUrl : `https://${jsUrl}`;
        const jsContent = await fetchUrl(fullUrl);
        
        if (jsContent.includes('WAWebVoipWebWasmLoader') || jsContent.includes('createVoipModule')) {
          console.log('[WASM Loader] Found loader in:', fullUrl);
          return fullUrl;
        }
      }
    }
  } catch (err) {
    console.warn('[WASM Loader] Loader discovery failed:', err.message);
  }
  return null;
}

function fetchUrl(url) {
  return new Promise((resolve, reject) => {
    https.get(url, { headers: { 'User-Agent': 'Mozilla/5.0' } }, (res) => {
      let data = '';
      res.on('data', chunk => data += chunk);
      res.on('end', () => resolve(data));
    }).on('error', reject);
  });
}

function fetchHead(url) {
  return new Promise((resolve, reject) => {
    const req = https.request(url, { method: 'HEAD', headers: { 'User-Agent': 'Mozilla/5.0' } }, resolve);
    req.on('error', reject);
    req.end();
  });
}

function downloadFile(url, destPath) {
  return new Promise((resolve, reject) => {
    const file = require('fs').createWriteStream(destPath);
    https.get(url, { headers: { 'User-Agent': 'Mozilla/5.0' } }, (res) => {
      if (res.statusCode !== 200) {
        reject(new Error(`HTTP ${res.statusCode}: ${url}`));
        return;
      }
      res.pipe(file);
      file.on('finish', () => file.close(resolve));
      file.on('error', reject);
    }).on('error', reject);
  });
}

// CLI kullanımı
if (import.meta.url === `file://${process.argv[1]}`) {
  const args = process.argv.slice(2);
  const urlArg = args.find(a => a.startsWith('--url='));
  const manualUrl = urlArg ? urlArg.split('=')[1] : null;
  
  (async () => {
    try {
      if (manualUrl) {
        const wasmPath = join(CACHE_DIR, `wa-voip-manual-${Date.now()}.wasm`);
        await downloadFile(manualUrl, wasmPath);
        console.log('Downloaded to:', wasmPath);
      } else {
        const path = await fetchLatestWasm();
        console.log('WASM ready at:', path);
      }
    } catch (err) {
      console.error('Failed:', err.message);
      process.exit(1);
    }
  })();
}