// !谨慎使用...解构，性能较差
// 解构: 140.362ms
// 普通赋值: 7.647ms

const obj = { a: 1, b: '2', c: true, d: 11 }

console.time('解构')
for (let i = 0; i < 1e7; i++) {
  const copy = { ...obj }
}
console.timeEnd('解构')

console.time('普通赋值')
for (let i = 0; i < 1e7; i++) {
  const copy = { a: obj.a, b: obj.b, c: obj.c, d: obj.d }
}
console.timeEnd('普通赋值')
