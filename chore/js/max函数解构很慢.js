// !max(...nums)还有一个缺点是,当参数大于1.25e5时,会爆栈(RangeError: Maximum call stack size exceeded)

const arr = Array(1e5) // 1.26e5
  .fill(0)
  .map((v, i) => i)

console.time('max...')
for (let i = 0; i < 1000; i++) {
  Math.max(...arr)
}
console.timeEnd('max...') // 290ms

console.time('naiveMax...')
for (let i = 0; i < 1000; i++) {
  let max = 0
  for (let j = 0; j < arr.length; j++) {
    max = Math.max(max, arr[j])
  }
}
console.timeEnd('naiveMax...') // 60ms

console.time('dont use max')
for (let i = 0; i < 1000; i++) {
  let max = 0
  for (let j = 0; j < arr.length; j++) {
    if (arr[j] > max) {
      max = arr[j]
    }
  }
}
console.timeEnd('dont use max') // 60ms
