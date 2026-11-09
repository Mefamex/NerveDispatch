const https = require('https');
const fs = require('fs');

https.get('https://static.whatsapp.net/rsrc.php/', (res) => {
  let data = '';
  res.on('data', chunk => data += chunk);
  res.on('end', () => {
    fs.writeFileSync('rsrc-dump.txt', data);
    console.log(`Saved ${data.length} chars to rsrc-dump.txt`);
    
    // .wasm ara
    const wasmMatches = data.match(/[a-zA-Z0-9_-]+\.wasm/g);
    if (wasmMatches) {
      console.log('\n=== .wasm filenames ===');
      console.log([...new Set(wasmMatches)].slice(0, 50));
    }
    
    // /rsrc.php/.../...wasm pattern
    const rsrcWasm = data.match(/rsrc\.php\/[^"'\s]*\.wasm/g);
    if (rsrcWasm) {
      console.log('\n=== rsrc.php wasm paths ===');
      console.log([...new Set(rsrcWasm)].slice(0, 50));
    }
  });
}).on('error', console.error);
