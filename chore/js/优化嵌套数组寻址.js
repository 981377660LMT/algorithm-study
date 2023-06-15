// 没用的优化

const n = 1e7
const a = Array(1e7)
  .fill(0)
  .map((_, i) => i)
const b = Array(1e7)
  .fill(0)
  .map((_, i) => i)
const c = Array(1e7)
  .fill(0)
  .map((_, i) => i)
const d = Array(1e7)
  .fill(0)
  .map((_, i) => i)

console.time('for')
for (let i = 0; i < n; i++) {
  a[i] = b[c[d[i]]]
}
console.timeEnd('for')

console.time('for2')
const helper = Array(n)
for (let i = 0; i < n; i++) {
  helper[i] = c[d[i]]
}
for (let i = 0; i < n; i++) {
  a[i] = b[helper[i]]
}
console.timeEnd('for2')
