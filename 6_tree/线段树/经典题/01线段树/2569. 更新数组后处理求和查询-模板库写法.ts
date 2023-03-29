import { useAtcoderLazySegmentTree } from '../../template/atcoder_segtree/AtcoderLazySegmentTree'

// 给你两个下标从 0 开始的数组 nums1 和 nums2 ，和一个二维数组 queries 表示一些操作。总共有 3 种类型的操作：
// 操作类型 1 为 queries[i] = [1, l, r] 。你需要将 nums1 从下标 l 到下标 r 的所有 0 反转成 1 或将 1 反转成 0 。l 和 r 下标都从 0 开始。
// 操作类型 2 为 queries[i] = [2, p, 0] 。对于 0 <= i < n 中的所有下标，令 nums2[i] = nums2[i] + nums1[i] * p 。
// 操作类型 3 为 queries[i] = [3, 0, 0] 。求 nums2 中所有元素的和。
// 请你返回一个数组，包含所有第三种操作类型的答案。
function handleQuery(nums1: number[], nums2: number[], queries: number[][]): number[] {
  const leaves = nums1.map(num => [num, 1])
  const seg01 = useAtcoderLazySegmentTree(leaves, {
    e: () => [0, 1],
    id: () => 0,
    op: ([sum1, size1], [sum2, size2]) => [sum1 + sum2, size1 + size2],
    mapping: (f, [sum, size]) => (f === 0 ? [sum, size] : [size - sum, size]),
    composition: (f, g) => f ^ g
  })

  let sum = nums2.reduce((a, b) => a + b, 0)
  const res: number[] = []
  queries.forEach(([op, a, b]) => {
    if (op === 1) {
      seg01.update(a, b + 1, 1)
    } else if (op === 2) {
      const ones = seg01.queryAll()[0]
      sum += ones * a
    } else {
      res.push(sum)
    }
  })

  return res
}
