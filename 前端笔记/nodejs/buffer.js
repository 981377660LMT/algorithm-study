const os = require('os')
const arr = new Uint16Array(2)
arr[0] = 5000
arr[1] = 4000

const buf1 = Buffer.from(arr) // 拷贝了该 buffer
const buf2 = Buffer.from(arr.buffer) // 与该数组共享了内存

console.log(buf1)
// 输出: <Buffer 88 a0>, 拷贝的 buffer 被截断只有两个元素
console.log(buf2) // 0x1388  0x0fa0
// 输出: <Buffer 88 13 a0 0f>

arr[1] = 6000
console.log(buf1)
// 输出: <Buffer 88 a0>
console.log(buf2)
// 输出: <Buffer 88 13 70 17>
console.log(os.EOL)

console.log(parseInt('0fa0', 16))
