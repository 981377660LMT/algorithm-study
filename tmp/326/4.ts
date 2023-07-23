import { SegmentTreeDynamic } from '../../6_tree/线段树/template/动态开点/SegmentTreeDynamicSparse'

export {}

const INF = 2e15

// 给你一个下标从 0 开始的整数数组 nums 和一个正整数 x 。

// 你 一开始 在数组的位置 0 处，你可以按照下述规则访问数组中的其他位置：

// 如果你当前在位置 i ，那么你可以移动到满足 i < j 的 任意 位置 j 。
// 对于你访问的位置 i ，你可以获得分数 nums[i] 。
// 如果你从位置 i 移动到位置 j 且 nums[i] 和 nums[j] 的 奇偶性 不同，那么你将失去分数 x 。
// 请你返回你能得到的 最大 得分之和。

// 注意 ，你一开始的分数为 nums[0] 。
function maxScore(nums: number[], x: number): number {
  const min = Math.min(...nums)
  const max = Math.max(...nums)
  const dp = [
    new SegmentTreeDynamic(min, max + 1, () => -INF, Math.max),
    new SegmentTreeDynamic(min, max + 1, () => -INF, Math.max)
  ]
  const n = nums.length
  dp[nums[0] & 1].update(nums[0], nums[0])
  for (let i = 1; i < n; ++i) {
    const num = nums[i]
    const odd = num & 1
    const cand1 = dp[odd].queryAll() + num
    const cand2 = dp[odd ^ 1].queryAll() + num - x
    dp[odd].update(num, Math.max(cand1, cand2))
  }
  return Math.max(dp[0].queryAll(), dp[1].queryAll())
}

// nums = [2,3,6,1,9,2], x = 5

console.log(maxScore([2, 3, 6, 1, 9, 2], 5))
