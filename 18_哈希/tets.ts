const aaaaasa = new Uint8Array(12)
console.log(aaaaasa.BYTES_PER_ELEMENT)

const buf = Buffer.alloc(5)
console.log(buf.byteLength)

// interface Buffer extends Uint8Array
// 注意Buffer继承自Uint8Array
