const nums1 = Array(1e5).fill(100000)
console.time('splice')
for (let i = 0; i < 5e4; i++) {
  nums1.splice(1000, 1)
  nums1.splice(1000, 1, 2)
}
console.timeEnd('splice') // splice: 3.239s

// !动态数组时 splice删除最快

let nums2 = new Uint32Array(1e5).map((_, i) => i)
console.time('copyWithin')
for (let i = 0; i < 1e5; i++) {
  // remove element at index 1000  使用copyWithin删除元素
  // nums2.copyWithin(1000, 1001) // copyWithin: 905.406ms
  nums2.set(nums2.subarray(1001), 1000) // set: 903.384ms
}
console.timeEnd('copyWithin')

// !静态数组时 copyWithin/subarray删除最快(但是不支持插入)
// !注意静态数组subarray比slice快
