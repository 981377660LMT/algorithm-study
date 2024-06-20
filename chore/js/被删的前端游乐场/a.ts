const order = new Map()
for (let i = 0; i < 200; i++) {
  order.set(i, ~~(Math.random() * 100))
}

const arr = Array.from({ length: 100 }, (_, i) => i + 1)

const f = new Map<number, number[]>()
for (let i = 0; i < 1e5; i++) {
  const arr = Array.from({ length: 20 }, () => ~~(Math.random() * 100))
  f.set(i, arr)
}

console.time('foo')
const res = []
for (let i = 0; i < 5e4; i++) {
  const arr = f.get(i)!
  arr.sort((a, b) => order.get(a) - order.get(b))
  f.set(i, arr)
}
console.timeEnd('foo')

export {}
