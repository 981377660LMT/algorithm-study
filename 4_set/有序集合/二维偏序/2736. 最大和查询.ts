// https://leetcode.cn/problems/maximum-sum-queries/
// 2736. 最大和查询 (二维偏序+离线查询)
// 对于第 i 个查询，在所有满足 nums1[j] >= xi 且 nums2[j] >= yi 的下标 j (0 <= j < n) 中，
// 找出 nums1[j] + nums2[j] 的 最大值 ，
// 如果不存在满足条件的 j 则返回 -1 。
// 返回数组 answer ，其中 answer[i] 是第 i 个查询的答案。
//
// !即:对每个查询(x,y),求出右上角的点的`横坐标+纵坐标`的最大值

import { SegmentTreeDynamic } from '../../../6_tree/线段树/template/动态开点/SegmentTreeDynamicSparse'

const INF = 2e15

function maximumSumQueries(nums1: number[], nums2: number[], queries: number[][]): number[] {
  const points = nums1.map((v, i) => [v, nums2[i]]).sort((a, b) => a[0] - b[0] || a[1] - b[1])
  const qWithId = queries.map((q, i) => [q[0], q[1], i]).sort((a, b) => a[0] - b[0] || a[1] - b[1])

  const seg = new SegmentTreeDynamic<number>(0, 1e9 + 10, () => -INF, Math.max)
  const res = Array(queries.length).fill(-1)
  let pi = points.length - 1
  for (let i = qWithId.length - 1; i >= 0; i--) {
    const [qx, qy, qid] = qWithId[i]
    while (pi >= 0 && points[pi][0] >= qx) {
      seg.update(points[pi][1]!, points[pi][0] + points[pi][1])
      pi--
    }
    const curMax = seg.query(qy!, 1e9 + 10)
    res[qid] = curMax === -INF ? -1 : curMax
  }

  return res
}
