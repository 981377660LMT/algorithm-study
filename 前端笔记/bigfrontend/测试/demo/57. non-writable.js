const a = {}
Object.defineProperty(a, 'foo1', {
  value: 1,
  // writable: true, // 默认false
})
const b = Object.create(a)
b.foo2 = 1

console.log(b.foo1)
console.log(b.foo2)

b.foo1 = 2
b.foo2 = 2
console.log(b)
console.log(b.foo1)
console.log(b.foo2)
// 1
// 1
// 1  因为foo1 不是 writable的 不能改
// 2
console.log(b.foo1)
