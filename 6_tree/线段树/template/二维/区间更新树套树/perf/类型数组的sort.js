// !大数组sort之前可以先转类型数组, sort之后再转回来(不要求原地排序时)

const arr = Array(1e7).fill(0 | (Math.random() * 1e9))

console.time('normal sort')
arr.sort((a, b) => a - b)
console.timeEnd('normal sort') // normal sort: 183.821ms

console.time('typed sort')
const arr2 = new Int32Array(arr)
arr2.sort() // 不要传compare，否则会慢很多
const normalArr = Array(arr2)
console.timeEnd('typed sort') // typed sort: 56.19ms

const foo = new Int32Array([11, 2, 12345, 3, 456, 0, 2])
foo.sort()
console.log(foo)
// Int32Array(7) [
//   0,  2,   2,
//   3, 11, 456,
// 12345
// ]
const bar = [11, 2, 12345, 3, 456, 0, 2]
bar.sort()
console.log(bar)
// [
//   0, 11, 12345,
//   2,  2,     3,
// 456
// ]
