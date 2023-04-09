// 0-9
const aaa = new Uint8Array(1e5).map(() => ~~(Math.random() * 10))
const bbb = new Uint8Array(1e5).map(() => ~~(Math.random() * 10))
console.time('a')
for (let i = 0; i < 1e2; i++) {
  aaa.fill(0, 1e5, ~~(Math.random() * 1e5))
  const s = aaa.subarray(0, ~~(Math.random() * 1e5)).toString()
  const s2 = bbb.subarray(0, ~~(Math.random() * 1e5)).toString()
  if (s === s2) {
  }
}
console.timeEnd('a')
