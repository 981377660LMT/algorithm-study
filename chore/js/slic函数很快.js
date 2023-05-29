const arr = Array(1e5)
  .fill(0)
  .map((v, i) => i)

console.time('slice...')
for (let i = 0; i < 1000; i++) {
  const copy = arr.slice(0)
}
console.timeEnd('slice...') // 370ms

console.time('for')
for (let i = 0; i < 1000; i++) {
  const copy = Array(arr.length)
  for (let j = 0; j < arr.length; j++) {
    copy[j] = arr[j]
  }
}
console.timeEnd('for') // 540ms
