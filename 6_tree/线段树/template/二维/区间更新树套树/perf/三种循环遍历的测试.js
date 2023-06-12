// 三种循环遍历的测试
// https://measurethat.net/Benchmarks/Show/2090/0/array-loop-vs-foreach-vs-map#latest_results_block
// forE: 118.373ms
// !for: 14.127ms
// map: 180.932ms

var arr = []
for (var i = 0; i < 1e7; i++) {
  arr[i] = i
}

function someFn(i) {
  return i * 3 * 8
}

console.time('forE')
arr.forEach(function (item) {
  someFn(item)
})
console.timeEnd('forE')

console.time('for')
for (var i = 0, len = arr.length; i < len; i++) {
  someFn(arr[i])
}
console.timeEnd('for')

console.time('map')
arr.map(item => someFn(item))

console.timeEnd('map')
