const lazy = Array(1e7).fill(1)

// without cache
console.time('without cache')
for (let i = 0; i < 1e7; ++i) {
  const a = lazy[i]
  const b = lazy[i]
  const c = lazy[i]
  const d = lazy[i]
  const e = lazy[i]
}
console.timeEnd('without cache') // without cache: 8.334ms

// with cache
console.time('with cache')
for (let i = 0; i < 1e7; ++i) {
  const cache = lazy[i]
  const a = cache
  const b = cache
  const c = cache
  const d = cache
  const e = cache
}
console.timeEnd('with cache') // with cache: 5.244ms
