const arr = Array(2000)
  .fill(0)
  .map((item, index) => index)

console.time('splice')
for (let i = 0; i < 2e5; i++) {
  arr.splice(0, 1)
  arr.push(2000 + i)
}
console.timeEnd('splice')
