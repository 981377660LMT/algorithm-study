// reverse很慢，for很快,且一个变量更快
// reverse: 469.526ms
// !for: 12.771ms
// for2: 14.188ms

const arr = Array(1e7).fill(0)

console.time('reverse')

arr.reverse()
console.timeEnd('reverse')

console.time('for')
for (let i = 0; i < arr.length / 2; i++) {
  const temp = arr[i]
  arr[i] = arr[arr.length - i - 1]
  arr[arr.length - i - 1] = temp
}
console.timeEnd('for')

console.time('for2')
for (let i = 0, j = arr.length - 1; i < j; i++, j--) {
  const temp = arr[i]
  arr[i] = arr[j]
  arr[j] = temp
}
console.timeEnd('for2')
