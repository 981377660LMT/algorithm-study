// for: 7.261ms
// for2: 4.19ms
// !手动循环展开

function bar(v) {}

console.time('for')
for (let i = 0; i < 1e7; i++) {
  bar(i)
}
console.timeEnd('for')

console.time('for2')
for (let i = 0; i < 1e7; ) {
  bar(i)
  i++
  bar(i)
  i++
  bar(i)
  i++
  bar(i)
  i++
}
console.timeEnd('for2')

const a = Array(100 + 4).fill(0)
let sum0 = 0
let sum1 = 0
let sum2 = 0
let sum3 = 0
for (let i = 0; i < 100; i += 4) {
  sum0 += a[i]
  sum1 += a[i + 1]
  sum2 += a[i + 2]
  sum3 += a[i + 3]
}
const sum = sum0 + sum1 + sum2 + sum3
