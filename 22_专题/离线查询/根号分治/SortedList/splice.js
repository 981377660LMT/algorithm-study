// splice in array

const arr = Array(1000)
  .fill(0)
  .map((_, i) => i)

const n = 2e5
const randomHead = Array(n)
for (let i = 0; i < n; i++) {
  randomHead[i] = Math.floor(10 * Math.random())
}

const time1 = performance.now()
for (let i = 0; i < n; i++) {
  arr.splice(randomHead[i], 1)
  arr.splice(randomHead[i], 0, i)
}
const time2 = performance.now()

console.log(time2 - time1) // 232.5
