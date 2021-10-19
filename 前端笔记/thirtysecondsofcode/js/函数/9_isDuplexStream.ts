import { Duplex, Stream } from 'stream'
console.log(isDuplexStream(new Stream.Duplex())) // true)

function isDuplexStream(val: any): val is Duplex {
  return (
    val !== null &&
    typeof val === 'object' &&
    typeof val.pipe === 'function' &&
    typeof val._read === 'function' &&
    typeof val._readableState === 'object' &&
    typeof val._write === 'function' &&
    typeof val._writableState === 'object'
  )
}
