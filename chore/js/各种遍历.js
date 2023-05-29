const arr = Array(1e5)
  .fill(0)
  .map((v, i) => i)

console.time('for of')
for (let i = 0; i < 1000; i++) {
  for (const v of arr) {
    v
  }
}
console.timeEnd('for of') // 65ms

console.time('for')
for (let i = 0; i < 1000; i++) {
  for (let j = 0; j < arr.length; j++) {
    arr[j] = 1
  }
}
console.timeEnd('for') // 30ms

console.time('for Each')
for (let i = 0; i < 1000; i++) {
  arr.forEach(v => {})
}
console.timeEnd('for Each') // !1s
