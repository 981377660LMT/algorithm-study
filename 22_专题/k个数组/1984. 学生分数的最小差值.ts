// 从数组中选出任意 k 名学生的分数，使这 k 个分数间 最高分 和 最低分 的 差值 达到 最小化 。
function minimumDifference(nums: number[], k: number): number {
  nums.sort((a, b) => a - b)
  let res = Infinity

  for (let i = k - 1; i < nums.length; i++) {
    res = Math.min(res, nums[i] - nums[i - k + 1])
  }

  return res
}

console.log(minimumDifference([9, 4, 1, 7], 2))

export default 1
