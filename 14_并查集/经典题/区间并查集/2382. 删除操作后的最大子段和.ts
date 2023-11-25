// 2382. 删除操作后的最大子段和
// https://leetcode.cn/problems/maximum-segment-sum-after-removals/

import { WeightedUnionFind } from '../../UnionFindWeighted-分量和'

function maximumSegmentSum(nums: number[], removeQueries: number[]): number[] {
  const n = nums.length
  const q = removeQueries.length

  const uf = new WeightedUnionFind(n)
  const res: number[] = Array(q).fill(0)
  let maxPartSum = 0
  const visited = new Uint8Array(n)
  for (let qi = q - 1; ~qi; qi--) {
    res[qi] = maxPartSum
    const pos = removeQueries[qi]
    uf.addWeight(pos, nums[pos])
    visited[pos] = 1
    if (pos > 0 && visited[pos - 1]) uf.union(pos, pos - 1)
    if (pos < n - 1 && visited[pos + 1]) uf.union(pos, pos + 1)
    maxPartSum = Math.max(maxPartSum, uf.getGroupWeight(pos))
  }

  return res
}
