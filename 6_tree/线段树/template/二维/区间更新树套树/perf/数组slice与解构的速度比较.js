// https://javascript.plainenglish.io/es5-vs-es6-performance-comparisons-c3606a241633

// !浅拷贝数组,slice更快(slice切片方法经过了优化,而解构本质上是for of,自然很慢)

const arr = Array(1e7).fill(0 | (Math.random() * 1e9))

console.time('slice')
const arr2 = arr.slice()
console.timeEnd('slice') // slice: 41.487ms

console.time('destructure')
const arr3 = [...arr]
console.timeEnd('destructure') // destructure: 48.672ms
