/**
 * @param {number[]} nums
 * @return {number}
 * 请你返回使 nums 严格递增 的 最少 操作次数。
 * 限制：只能对某个数进行+1
 * 如果可以加一也可以减一呢?
 * 用 LIS 就可以解答任意调整的情况
 */
function minOperations(nums) {
  const n = nums.length
  if (n <= 1) return 0

  let preMax = nums[0]
  let res = 0

  for (let i = 1; i < n; i++) {
    if (preMax < nums[i]) {
      preMax = nums[i]
    } else {
      const diff = preMax - nums[i]
      res += diff + 1
      preMax += 1
    }
  }

  return res
}

// 输入：nums = [1,1,1]
// 输出：3
// 解释：你可以进行如下操作：
// 1) 增加 nums[2] ，数组变为 [1,1,2] 。
// 2) 增加 nums[1] ，数组变为 [1,2,2] 。
// 3) 增加 nums[2] ，数组变为 [1,2,3] 。
