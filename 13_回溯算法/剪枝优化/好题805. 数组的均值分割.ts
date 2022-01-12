const EPS = 1e-6
function splitArraySameAverage(nums: number[]): boolean {
  if (nums.length <= 1) return false
  const n = nums.length
  const avg = nums.reduce((pre, cur) => pre + cur, 0) / n

  for (let i = 1; i <= nums.length >> 1; i++) {
    // 剪枝2 和必须是整数 取模检查小数部分
    if (Math.abs((avg * i) % 1) > EPS) continue
    if (checkSum(nums, i, avg * i)) return true
  }

  return false

  /**
   *
   * @param nums
   * @param i 选择i个数的和为target
   */
  function checkSum(nums: number[], i: number, target: number): boolean {
    nums.sort((a, b) => a - b)

    const bt = (index: number, pathSum: number, remain: number, target: number): boolean => {
      // 剪枝3
      if (pathSum > target) return false

      if (remain === 0) {
        return Math.abs(pathSum - target) < EPS
      }

      for (let i = index; i < nums.length; i++) {
        // 剪枝1
        if (i !== index && nums[i] === nums[i - 1]) continue
        if (bt(i + 1, pathSum + nums[i], remain - 1, target)) return true
      }

      return false
    }

    return bt(0, 0, i, target)
  }
}

// 划分Nums 使两个子数组平均值相等 (都等于元数组平均值)
console.log(splitArraySameAverage([1, 2, 3, 4, 5, 6, 7, 8]))
console.log(splitArraySameAverage([2, 12, 18, 16, 19, 3]))

// 记录：
// 50 / 92 个通过测试用例
// FATAL ERROR: CALL_AND_RETRY_LAST Allocation failed - JavaScript heap out of memory
// [60,30,30,30,30,30,30,30,30,30,30,30,30,30,30,30,30,30,30,30,30,30,30,30,30,30,30,30,30,30]
// 1.不看相同的元素 sort + 相等时只保留第一个

// 68 / 92 个通过测试用例
// FATAL ERROR: CALL_AND_RETRY_LAST Allocation failed - JavaScript heap out of memory
// [33,86,88,78,21,76,19,20,88,76,10,25,37,97,58,89,65,59,98,57,50,30,58,5,61,72,23,6]
// 2.由于每个数字都是整数，那么其和一定也是整数
