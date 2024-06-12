// https://chinese.freecodecamp.org/forum/t/topic/633
// !1.不要用Array.from (es6新增语法都比较慢)
// !2.Array(n).fill(0) 比 Array.fill(undefined) 等别的填充物快
// !只要初始化时,甚至不需要fill,直接用Array(n)就可以了
// 在工作中不要为了新语法而ES6/7新语法，灵活使用，前端也是要注意性能和算法，
// 老的api在这方面应该是有明显的长处的
console.time('Array.from')
const arr2 = Array.from({ length: 1e7 }, () => [])
console.timeEnd('Array.from') // Array.from: 1.344s

console.time('array')
const arr = Array(1e7)
  .fill(0)
  .map(() => [])
console.timeEnd('array') // Array: 426.143ms

console.time('map')
const arr3 = Array(1e7)
  .fill(0)
  .map((_, i) => i)
console.timeEnd('map') // 232.915ms

console.time('for')
const arr4 = Array(1e7).fill(0)
for (let i = 0; i < arr4.length; i++) {
  arr4[i] = i
}
console.timeEnd('for') // 64.761ms

export {}

// !优化算法时，只用es5语法最快
// 特例： ForEach 循环很快
