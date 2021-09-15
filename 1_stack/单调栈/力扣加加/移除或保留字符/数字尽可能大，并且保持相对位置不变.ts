/**
 *
 * @param nums1
 * @param nums2
 * @description
 * 数字尽可能大，并且保持相对位置不变
 */
const pickMax = (nums1: number[], nums2: number[]): number[] => {
  // 一直shift 先转成deque比较好
  const res: number[] = []

  while (nums1.length || nums2.length) {
    // 数组比大小隐式转换成数字(parseInt)
    const bigger = nums1 > nums2 ? nums1 : nums2
    res.push(bigger.shift()!)
  }
  return res
}

console.log(pickMax([4, 4, 1], [5, 4, 3]))
console.log(pickMax([6, 7], [6, 0, 4]))
