const set = new Set<any>()

set.add(1).add(5).add(5).add({ a: 1, b: 2 }).add({ a: 1, b: 2 }).delete(5)

const has = set.has(5)

console.log(set.size)

// 迭代set
// for (const item of set) {
//   console.log(item)
// }
// for (const item of set.keys()) {
//   console.log(item)
// }
// for (const item of set.values()) {
//   console.log(item)
// }
for (const item of set.entries()) {
  console.log(item)
}
export {}

// key和value一摸一样

// set 转array
// const arr = [...set]
// const arr = Array.from(set)
// console.log(arr)

// array转set
const set2 = new Set([1, 2, 3, 4])

// 求set交集
const intersection = new Set([...set].filter(ele => set2.has(ele)))
console.log(intersection)
