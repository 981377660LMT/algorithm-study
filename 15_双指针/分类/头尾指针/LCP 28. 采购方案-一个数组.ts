const MOD = 1e9 + 7

/**
 * @param nums  2 <= nums.length <= 10^5
 * 假定小力仅购买两个零件，要求购买零件的花费不超过预算，请问他有多少种采购方案
 * !注意(a,b)和(b,a)是同一种方案
 */
function purchasePlans(nums: number[], target: number): number {
  const n = nums.length
  nums.sort((a, b) => a - b)

  let res = 0
  let left = 0
  let right = n - 1

  while (left < right) {
    if (nums[left] + nums[right] <= target) {
      res += right - left
      res %= MOD
      left++
    } else {
      right--
    }
  }

  return res
}

console.log(purchasePlans([2, 2, 1, 9], 10))
// 小力预算为 target，假定小力仅购买两个零件，要求购买零件的花费不超过预算，请问他有多少种采购方案。

export {}
