/**
 * @param {number[]} nums  1 <= nums[i] <= 1000
 * @param {number} k
 * @return {number}
 * 请找出该数组内乘积小于 k 的连续的子数组的个数。
 * @summary
 * 窗口内每加入一个数，如果乘积仍然小于k，则 满足条件的子数组数量 += 窗口内元素数量
 * 当遇到比k还大的元素时，直接重新在这个元素的右边继续开始计算
   当窗口内元素乘积大于等于k时，在左边踢出元素，收缩窗口，直到乘积不再大于等于k
 */
const numSubarrayProductLessThanK = function (nums: number[], k: number): number {
  if (k === 0) return 0

  let res = 0
  let left = 0
  let right = 0
  let product = 1

  while (right < nums.length) {
    const cur = nums[right]
    if (cur >= k) {
      left = right = right + 1
      product = 1
      continue
    }

    product = product * cur
    while (product >= k) {
      product /= nums[left]
      left++
    }

    res += right - left + 1
    right++
  }

  return res
}

console.log(numSubarrayProductLessThanK([10, 5, 2, 6], 100))
// 8
// 解释: 8个乘积小于100的子数组分别为: [10], [5], [2], [6], [10,5], [5,2], [2,6], [5,2,6]。

export {}
