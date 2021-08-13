/**
 * @param {number[]} nums1
 * @param {number[]} nums2
 * @param {number} k
 * @return {number[]}
 * 现在从这两个数组中选出 k (k <= m + n) 个数字拼接成一个新的数，要求从同一个数组中取出的数字保持其在原数组中的相对顺序。
 */
const maxNumber = (nums1: number[], nums2: number[], k: number): number[] => {
  let max: number[] = []
  /**
   * @description  见402. 移掉 K 位数字
   * @param arr 待选数组
   * @param k 顺序选出K个数保证组成的数字最大
   */
  const pickMaxKFromArray = (arr: number[], k: number) => {
    let drop = arr.length - k
    const stack: number[] = []
    for (let i = 0; i < arr.length; i++) {
      while (drop && stack.length && arr[stack[stack.length - 1]] < arr[i]) {
        stack.pop()
        drop--
      }
      stack.push(i)
    }
    return stack.map(index => arr[index])
  }

  /**
   *
   * @param arr1 归并排序的数组
   * @param arr2 归并排序的数组
   */
  const merge = (nums1: number[], nums2: number[]) => {
    const arr1 = nums1.slice()
    const arr2 = nums2.slice()
    const res: number[] = []

    while (arr1.length && arr2.length) {
      if (arr1[0] < arr2[0]) res.push(arr1.shift()!)
      else res.push(arr2.shift()!)
    }

    return [...res, ...arr1, ...arr2]
  }

  for (let i = 0; i <= k; i++) {
    if (i <= nums1.length && k - i <= nums2.length) {
      const tmp = merge(pickMaxKFromArray(nums1, i), pickMaxKFromArray(nums2, k - i))
      // 数组比大小隐式转换成数字(parseInt)
      if (tmp > max) max = tmp
    }
  }

  return max
}

console.log(maxNumber([3, 4, 6, 5], [9, 1, 2, 5, 8, 3], 5))
