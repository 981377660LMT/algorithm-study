// The copyWithin() method shallow copies part of an array to another location
// in the same array and returns it without modifying its length.

// The copyWithin works like C and C++'s memmove,
// and is a high-performance method to shift the data of an Array
// copyWithin 就像 C 和 C++ 的 memcpy 函数一样，且它是用来移动 Array 或者 TypedArray 数据的一个高性能的方法。
// 复制以及粘贴序列这两者是为一体的操作;即使复制和粘贴区域重叠，粘贴的序列也会有拷贝来的值。

// python切片/js的copywithin都是调用C的memmove
// start和end都是可以省略。
// start省略表示从0开始，end省略表示数组的长度值。

const arr = new Uint8Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
console.log(arr.subarray(0, 3)) // !快速切片
console.log(arr.copyWithin(0, -3, -1)) // !从(-3,-1)切片开始拷贝 从0开始覆盖
arr.set([1, 2, 3, 4, 9, 9, 9, 9, 9], 1) // !切片赋值 注意不能超出数组长度
// !切片赋值 注意不能超出数组长度()
// arr.set([1, 2, 3, 4, 9, 9, 9, 9, 9, 9], 1)
console.log(arr, 99)

// 把 nums1 序列中从下标 i1 位置开始的连续 k 个元素粘贴到 nums2 序列中从下标 i2 开始的连续 k 个位置上.
// 超出部分可忽略.
function copy(nums1: Int32Array, nums2: Int32Array, i1: number, k: number, i2: number): void {
  const len = Math.min(k, nums1.length - i1, nums2.length - i2)
  nums2.set(nums1.subarray(i1, i1 + len), i2)
}

if (require.main === module) {
  // test copy speed
  const arr1 = new Int32Array(1e5).fill(~~(Math.random() * 1e5))
  const arr2 = new Int32Array(1e5).fill(~~(Math.random() * 1e5))
  console.time('copy')
  for (let i = 0; i < 1e5; i++) {
    copy(arr1, arr2, 0, 1e5, 0)
  }
  console.timeEnd('copy') // 1.3s左右
}

export {}
