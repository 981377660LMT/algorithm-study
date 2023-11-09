// !遍历二维数组,缓存一下会稍快

const big2d = new Array(5000).fill(0).map(() => new Array(5000).fill(0))

console.time('dont cache')
for (let i = 0; i < 5000; ++i) {
  for (let j = 0; j < 5000; ++j) {
    const cur = big2d[i][j]
  }
}
console.timeEnd('dont cache')

console.time('cache')
for (let i = 0; i < 5000; ++i) {
  const row = big2d[i]
  for (let j = 0; j < 5000; ++j) {
    const cur = row[j]
  }
}
console.timeEnd('cache')
