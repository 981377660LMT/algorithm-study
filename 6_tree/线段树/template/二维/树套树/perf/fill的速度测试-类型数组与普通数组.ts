// fill的速度测试-类型数组与普通数组
//
// 类型数组：fill很快.
// 普通数组: fill基本类型很快,引用类型很慢(一般需要map等函数).
//
// !初始化时可以不fill，但是必须保证数组随机访问时不出现hole(顺序给数组下标赋值).
// 否则会倍永久标记成为稀疏数组。
// !并且fill时不要都fill成undefined,最好保持类型一样,即fill成待保存数据的类型.
// 有利于V8识别数组类型。
//
// ---------------------------------------------------
// 一个原则: 缓存友好性 => 根号算法用于js的优化
// 两个拥抱: 拥抱类型数组与原生JS
// 三个函数: fill, slice, splice (es6里优化过的)
// ---------------------------------------------------

const n = 1e7

console.time('fill')
const arr1 = Array(n).fill(123)
console.timeEnd('fill')

console.time('for')
const arr2 = Array(n)
for (let i = 0; i < n; i++) arr2[i] = 123
console.timeEnd('for')

console.time('fill')
const arr3 = new Uint32Array(n).fill(100)
console.timeEnd('fill')

console.time('for')
const arr4 = new Uint32Array(n)
for (let i = 0; i < n; i++) arr4[i] = 100
console.timeEnd('for')
export {}
