let arr = Array.from({ length: 100 }, (_, i) => i)
const set = new Set(arr)
console.time('foo')

for (let i = 0; i < 1e4; i++) {
  arr = arr.filter(v => set.has(v))
}
console.timeEnd('foo')

export {}
