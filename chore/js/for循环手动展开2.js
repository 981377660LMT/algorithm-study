// 求数组和
// 可能会刺激 CPU 并发运算，从而提高运行效率。

const rands = Array.from({ length: 1e7 }, () => (Math.random() * 1e8) | 0)

console.time('for')
let sum = 0
for (let i = 0; i < rands.length; i++) {
  sum += rands[i]
}
console.timeEnd('for')

console.time('for2')

let n = rands.length
let sum0 = 0
let sum1 = 0
let sum2 = 0
let sum3 = 0

switch (n & 3) {
  case 3:
    sum0 += rands[n - 3] + rands[n - 2] + rands[n - 1]
    break
  case 2:
    sum0 += rands[n - 2] + rands[n - 1]
    break
  case 1:
    sum0 += rands[n - 1]
    break
  default:
    break
}

n -= 3
for (let i = 0; i < n; i += 4) {
  sum0 += rands[i]
  sum1 += rands[i + 1]
  sum2 += rands[i + 2]
  sum3 += rands[i + 3]
}

console.timeEnd('for2')
