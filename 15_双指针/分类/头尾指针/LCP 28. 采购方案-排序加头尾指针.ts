/**
 *
 * @param nums  2 <= nums.length <= 10^5
 * @param target  假定小力仅购买两个零件，要求购买零件的花费不超过预算，请问他有多少种采购方案
 */
function purchasePlans(nums: number[], target: number): number {
  nums.sort((a, b) => a - b)
  let res = 0
  let l = 0
  let r = nums.length - 1
  while (l < r) {
    const sum = nums[l] + nums[r]
    if (sum <= target) {
      res += r - l
      l++
    } else r--
  }

  return res % (10 ** 9 + 7)
}

console.log(purchasePlans([2, 2, 1, 9], 10))
// 小力预算为 target，假定小力仅购买两个零件，要求购买零件的花费不超过预算，请问他有多少种采购方案。
