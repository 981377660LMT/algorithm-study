type SumMod = number
type Index = number

/**
 * @param {number[]} nums  1 <= nums.length <= 105
 * @param {number} k
 * @return {boolean}
 * 编写一个函数来判断该数组是否含有同时满足下述条件的连续子数组：
   子数组大小 至少为 2 ，且
   子数组元素总和为 k 的倍数。
   当两个不同位置的前缀和对 kk 的取余相同时，我们看这两个位置的下标是否距离大于等于2.
 */
const checkSubarraySum = function (nums: number[], k: number): boolean {
  // -1 位置和为0
  const first = new Map<SumMod, Index>([[0, -1]])
  let sum = 0

  for (let i = 0; i < nums.length; i++) {
    sum += nums[i]
    const mod = sum % k
    if (first.has(mod)) {
      if (i - first.get(mod)! >= 2) {
        return true
      }
    } else {
      first.set(mod, i)
    }
  }

  return false
}

console.log(checkSubarraySum([23, 2, 4, 6, 7], 6))
console.log(checkSubarraySum([23, 2, 4, 6, 6], 7))

export default 1
console.log(2 % 0) // NaN
