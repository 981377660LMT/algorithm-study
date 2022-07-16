const buffer = new ArrayBuffer(1024)
const arr = new Uint8Array(buffer)
const arr2 = new Uint32Array(buffer)
console.log(arr.length)
console.log(arr2.length)
arr2[0] = 1
console.log(arr)
type MutableMethod = Exclude<keyof Array<any>, keyof Uint8Array>

const arrayLike = { length: 2, a: 12, 0: 9 }
console.log(Array.from(arrayLike))
'as'.trimStart()

export {}
