// console.log(Array.prototype.reverse.call({ 1: 'banana', 2: 'apple', 3: 'orange', length: 4 }))

// 翻转对象键值
console.log(
  Object.fromEntries(
    Object.entries({ 1: 'banana', 2: 'apple', 3: 'orange' }).map(([a, b]) => [b, a])
  )
)

const map = new Map([
  ['k1', 'v1'],
  ['k2', 'v2'],
])

// 灵活使用[Symbol.iterable]统一的迭代接口 of
for (const [k, v] of map) {
  console.log(k, v)
}

export {}
