/**
 * @param {number[]} nums
 * @return {boolean}
 * @description 判断是否可以将这个数组分割成两个子集，使得两个子集的元素和相等
 * @summary 01背包问题:在n个物品中选一定物品完全填满sum/2容量的背包
 */
const canPartition = function (nums: number[]) {
  const sum = nums.reduce((pre, cur) => pre + cur, 0)
  const volume = sum / 2
  if (!Number.isInteger(volume)) return false
  const dpRow = nums.length
  const dp = Array<boolean>(volume + 1).fill(false)

  dp[0] = true

  for (let i = 0; i < dpRow; i++) {
    const num = nums[i]
    // 倒叙的原因是后面需要前面先决定
    for (let j = volume; j >= num; j--) {
      // dp[j]不使用第i个物品
      // dp[j-num]使用第i个物品
      dp[j] = dp[j] || dp[j - num]
    }
  }

  return dp[volume]
}

console.dir(canPartition([1, 2, 2, 5]), { depth: null })

export {}
