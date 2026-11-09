import https from 'https';
import { promises as fs, createWriteStream } from 'fs';
import { join } from 'path';

const CACHE_DIR = join(process.env.HOME || '/tmp', '.wa-voip-sidecar', 'cache');

async function fetchUrl(url) {
  return new Promise((resolve, reject) => {
    https.get(url, { headers: { 'User-Agent': 'Mozilla/5.0' } }, (res) => {
      let data = '';
      res.on('data', chunk => data += chunk);
      res.on('end', () => resolve(data));
    }).on('error', reject);
  });
}

async function discoverWasmUrl() {
  try {
    console.log('[Discovery] Fetching WhatsApp Web...');
    const html = await fetchUrl('https://web.whatsapp.com/');
    
    // JS bundle URL'lerini bul
    const jsUrls = html.match(/static\.whatsapp\.net\/rsrc\.php\/v4[a-zA-Z0-9/_-]+\.js/g);
    if (jsUrls) {
      const uniqueJsUrls = [...new Set(jsUrls)];
      console.log(`[Discovery] Found ${uniqueJsUrls.length} JS bundles, checking first 10...`);
      
      for (const jsUrl of uniqueJsUrls.slice(0, 10)) {
        const fullUrl = jsUrl.startsWith('http') ? jsUrl : `https://${jsUrl}`;
        try {
          const jsContent = await fetchUrl(fullUrl);
          
          // WASM referansı ara
          const wasmMatch = jsContent.match(/https?:\/\/static\.whatsapp\.net\/rsrc\.php\/[^"'\s]+\.wasm/);
          if (wasmMatch) {
            console.log('[Discovery] Found WASM URL:', wasmMatch[0]);
            return wasmMatch[0];
          }
          
          // Loader script referansı ara
          const loaderMatch = jsContent.match(/WAWebVoipWebWasmLoader[^"'\s]*/);
          if (loaderMatch) {
            console.log('[Discovery] Found loader reference:', loaderMatch[0]);
          }
        } catch (err) {
          // Ignore individual bundle errors
        }
      }
    }
  } catch (err) {
    console.warn('[Discovery] Auto-discovery failed:', err.message);
  }
  
  // Bilinen hash'leri dene
  const KNOWN_WASM_HASHES = ['KwJJIha0f3H'];
  for (const hash of KNOWN_WASM_HASHES) {
    const testUrl = `https://static.whatsapp.net/rsrc.php/ys/r/${hash}.wasm`;
    try {
      const res = await fetchHead(testUrl);
      if (res.ok) {
        console.log('[Discovery] Found working known hash:', testUrl);
        return testUrl;
      }
    } catch {}
  }
  
  return null;
}

async function discoverLoaderUrl() {
  try {
    const html = await fetchUrl('https://web.whatsapp.com/');
    const jsUrls = html.match(/static\.whatsapp\.net\/rsrc\.php\/v4[a-zA-Z0-9/_-]+\.js/g);
    
    if (jsUrls) {
      for (const jsUrl of [...new Set(jsUrls)].slice(0, 10)) {
        const fullUrl = jsUrl.startsWith('http') ? jsUrl : `https://${jsUrl}`;
        const jsContent = await fetchUrl(fullUrl);
        
        if (jsContent.includes('WAWebVoipWebWasmLoader') || jsContent.includes('createVoipModule')) {
          console.log('[Discovery] Found loader in:', fullUrl);
          return fullUrl;
        }
      }
    }
  } catch (err) {
    console.warn('[Discovery] Loader discovery failed:', err.message);
  }
  return null;
}

function fetchHead(url) {
  return new Promise((resolve, reject) => {
    const req = https.request(url, { method: 'HEAD', headers: { 'User-Agent': 'Mozilla/5.0' } }, resolve);
    req.on('error', reject);
    req.end();
  });
}

async function downloadFile(url, destPath) {
  return new Promise((resolve, reject) => {
    const file = createWriteStream(destPath);
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

async function main() {
  const args = process.argv.slice(2);
  const urlArg = args.find(a => a.startsWith('--url='));
  const manualUrl = urlArg ? urlArg.split('=')[1] : null;
  
  try {
    await fs.mkdir(CACHE_DIR, { recursive: true });
    
    if (manualUrl) {
      const wasmPath = join(CACHE_DIR, `wa-voip-manual-${Date.now()}.wasm`);
      await downloadFile(manualUrl, wasmPath);
      console.log('Downloaded to:', wasmPath);
    } else {
      const wasmUrl = await discoverWasmUrl();
      if (!wasmUrl) {
        console.error('Could not discover WASM URL. Try manually: node fetch-wasm.js --url <url>');
        process.exit(1);
      }
      
      const wasmPath = join(CACHE_DIR, `wa-voip-${Date.now()}.wasm`);
      await downloadFile(wasmUrl, wasmPath);
      
      const loaderUrl = await discoverLoaderUrl();
      if (loaderUrl) {
        const loaderPath = join(process.cwd(), 'WAWebVoipWebWasmLoader.js');
        await downloadFile(loaderUrl, loaderPath);
        console.log('Loader saved to:', loaderPath);
      }
      
      await fs.writeFile(
        join(CACHE_DIR, 'latest.json'),
        JSON.stringify({ wasmUrl, wasmPath, loaderUrl, timestamp: Date.now() }, null, 2)
      );
      
      console.log('WASM ready at:', wasmPath);
    }
  } catch (err) {
    console.error('Failed:', err.message);
    process.exit(1);
  }
}

main();