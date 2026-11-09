const https = require('https');

// WhatsApp Web'in manifest/asset JSON'ını çekmeyi dene
const urls = [
  'https://web.whatsapp.com/manifest.json',
  'https://static.whatsapp.net/rsrc.php/',
  'https://web.whatsapp.com/check-update?version=2.2400.0&platform=web',
];

async function fetchUrl(url) {
  return new Promise((resolve, reject) => {
    https.get(url, (res) => {
      let data = '';
      res.on('data', chunk => data += chunk);
      res.on('end', () => resolve({url, status: res.statusCode, data}));
    }).on('error', reject);
  });
}

(async () => {
  for (const url of urls) {
    try {
      const result = await fetchUrl(url);
      console.log(`\n=== ${url} (${result.status}) ===`);
      if (result.data.length < 5000) {
        console.log(result.data.substring(0, 2000));
      } else {
        console.log(`[${result.data.length} chars] searching for .wasm...`);
        const matches = result.data.match(/[a-zA-Z0-9_-]+\.wasm/g);
        if (matches) console.log([...new Set(matches)].slice(0, 20));
      }
    } catch (e) {
      console.log(`\n=== ${url} ERROR: ${e.message} ===`);
    }
  }
})();
