// !慢数组的负数下标访问比正数下标访问慢很多

const n = 1e6

const arr = []
console.time('a')
for (let i = 0; i < n; i++) {
  arr[-i] = i
  arr[-i]
}
console.timeEnd('a') // 839.943ms

const arr2 = []
console.time('b')
for (let i = 0; i < n; i++) {
  arr2[i] = i
  arr2[i]
}
console.timeEnd('b') // 20.373ms
