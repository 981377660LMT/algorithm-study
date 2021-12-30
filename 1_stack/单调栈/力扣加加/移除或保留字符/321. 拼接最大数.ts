/**
 * @param {number[]} nums1
 * @param {number[]} nums2
 * @param {number} k
 * @return {number[]}
 * 现在从这两个数组中选出 k (k <= m + n) 个数字拼接成一个新的数，要求从同一个数组中取出的数字保持其在原数组中的相对顺序。
 */
const maxNumber = (nums1: number[], nums2: number[], k: number): number[] => {
  let res: number[] = []

  for (let i = 0; i <= k; i++) {
    if (i <= nums1.length && k - i <= nums2.length) {
      const cand = merge(pickMaxK(nums1, i), pickMaxK(nums2, k - i))
      // 数组比大小隐式转换成数字(parseInt)
      if (cand > res) res = cand
    }
  }

  return res

  /**
   * @description  见402. 移掉 K 位数字
   * @param arr 待选数组
   * @param k 顺序选出K个数保证组成的数字最大
   */
  function pickMaxK(arr: number[], k: number) {
    let drop = arr.length - k
    const stack: number[] = []
    for (let i = 0; i < arr.length; i++) {
      while (drop > 0 && stack.length > 0 && stack[stack.length - 1] < arr[i]) {
        stack.pop()
        drop--
      }
      stack.push(arr[i])
    }
    // 注意slice
    return stack.slice(0, k)
  }

  /**
   *
   * @param arr1 归并排序的数组
   * @param arr2 归并排序的数组
   */
  function merge(nums1: number[], nums2: number[]) {
    // 一直shift 先转成deque比较好
    const res: number[] = []

    while (nums1.length || nums2.length) {
      // 数组比大小隐式转换成数字(parseInt)
      const bigger = nums1 > nums2 ? nums1 : nums2
      res.push(bigger.shift()!)
    }
    return res
  }
}

// console.log(maxNumber([3, 4, 6, 5], [9, 1, 2, 5, 8, 3], 5))
console.log(maxNumber([6, 7, 5], [4, 8, 1], 3))
