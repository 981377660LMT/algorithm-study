// 异或a ^ b ^ c = a ^ c ^ b
// 任何数和本身异或则为0 (可以检测在一串异或中是否重复出现过这个数字)
// 任何数和 0 异或是本身

import assert from 'assert'

console.log(1 ^ 3)
console.log(3 ^ 3)
console.log(5 ^ 0)
console.log((1 ^ 2 ^ 4) === (2 ^ 4 ^ 1))
console.log(1 ^ 2 ^ 1 ^ 3 ^ 3)
// 位运算优先级
console.log((14 >> 2) & 1)

// 注意：异或的MSB 最高有效位 (most significant bit)
// 由于异或中每个1都只来源于一个数，则MSB位的数必不相同，那么我们可以根据MSB的1将两个数区分开来
// 即向右位移MSB位，与1，一个必为1一个必为0
console.log(Number(3).toString(2))
console.log(Number(5).toString(2))
console.log(Number(3 ^ 5).toString(2))
const MSB = Number(3 ^ 5).toString(2).length - 1
assert.notStrictEqual((3 >> MSB) & 1, (5 >> MSB) & 1)
