const arw = { a: 1, b: 'qsq3434' }

console.time('push')
for (let i = 0; i < 1e7; i++) {
  const a = i * 2 + 1
}
console.timeEnd('push')

console.time('spread')
for (let i = 0; i < 1e7; i++) {
  const a = (i << 1) | 1
}
console.timeEnd('spread')
