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

// MSB启示我们只需要比较异或的前缀就可以看出一些性质
// 例如异或是否比某个数大
// 只需要比较前缀请即可
// 如果一个数 a 的前缀和另外一个数 b 的前缀是一样的，
// 那么 c 和 a 或者 c 和 b 的异或的结构前缀部分一定也是一样的

const foo = 0b110101

console.log(foo & -foo) // lowbit
console.log(foo & (foo - 1)) // 将最后一位 1 变成 0
console.log(foo | (foo + 1)) // 将最后一位 0 变成 1
console.log(foo & ((1 << 3) - 1)) // 将 foo 最高位至第 3 位(含)清零
console.log(foo & ~(1 << 3)) // 将第 n 位清零
