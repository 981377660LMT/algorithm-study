// reverse快
// !reverse:9.399999999906868
// for:14
// for2:11.70000000006985

let arr = Array(1e7).fill(0)
let time = performance.now()
arr.reverse()
console.log(performance.now() - time)

arr = Array(1e7).fill(0)
time = performance.now()
for (let i = 0; i < arr.length / 2; i++) {
  const temp = arr[i]
  arr[i] = arr[arr.length - i - 1]
  arr[arr.length - i - 1] = temp
}
console.log(performance.now() - time)

arr = Array(1e7).fill(0)
time = performance.now()
for (let i = 0, j = arr.length - 1; i < j; i++, j--) {
  const temp = arr[i]
  arr[i] = arr[j]
  arr[j] = temp
}
console.log(performance.now() - time)
