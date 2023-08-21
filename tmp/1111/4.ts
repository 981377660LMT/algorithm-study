import { SortedListFast } from '../../22_专题/离线查询/根号分治/SortedList/SortedListFast'

export {}

const INF = 2e15
// 给你一个下标从 0 开始长度为 n 的整数数组 nums 和一个整数 target ，请你返回满足 0 <= i < j < n 且 nums[i] + nums[j] < target 的下标对 (i, j) 的数目。

function countPairs(nums: number[], target: number): number {
  const sl = new SortedListFast<number>()
  let res = 0
  nums.forEach(v => {
    res += sl.bisectLeft(target - v)
    sl.add(v)
  })
  return res
}
