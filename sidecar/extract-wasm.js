const https = require('https');

https.get('https://static.whatsapp.net/rsrc.php/', (res) => {
  let data = '';
  res.on('data', chunk => data += chunk);
  res.on('end', () => {
    // .wasm dosyalarını bul
    const wasmUrls = data.match(/https?:\/\/[^"'\s]+\.wasm/g);
    const wasmFiles = data.match(/[a-zA-Z0-9_-]+\.wasm/g);
    const voipRefs = data.match(/voip[^"'\s]*/gi);
    const loaderRefs = data.match(/WAWebVoip[^"'\s]*/gi);
    
    console.log('=== Full .wasm URLs ===');
    console.log(wasmUrls ? [...new Set(wasmUrls)] : 'none');
    
    console.log('\n=== .wasm filenames ===');
    console.log(wasmFiles ? [...new Set(wasmFiles)].slice(0, 30) : 'none');
    
    console.log('\n=== voip refs ===');
    console.log(voipRefs ? [...new Set(voipRefs)].slice(0, 20) : 'none');
    
    console.log('\n=== WAWebVoip refs ===');
    console.log(loaderRefs ? [...new Set(loaderRefs)] : 'none');
    
    // Dosyaya kaydet
    fs.writeFileSync('rsrc-dump.txt', data);
    console.log('\nFull dump saved to rsrc-dump.txt');
  });
}).on('error', console.error);
