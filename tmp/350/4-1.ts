import { useAtcoderLazySegmentTree } from '../../6_tree/线段树/template/atcoder_segtree/AtcoderLazySegmentTree'

export {}

// 给你两个长度为 n 、下标从 0 开始的整数数组 nums1 和 nums2 ，另给你一个下标从 1 开始的二维数组 queries ，其中 queries[i] = [xi, yi] 。

// 对于第 i 个查询，在所有满足 nums1[j] >= xi 且 nums2[j] >= yi 的下标 j (0 <= j < n) 中，找出 nums1[j] + nums2[j] 的 最大值 ，如果不存在满足条件的 j 则返回 -1 。

// 返回数组 answer ，其中 answer[i] 是第 i 个查询的答案。

// !将查询排序,枚举答案
function maximumSumQueries(nums1: number[], nums2: number[], queries: number[][]): number[] {
  const res = Array(queries.length).fill(-1)
  const qWithId = queries.map((q, i) => [...q, i]).sort((a, b) => a[0] - b[0] || a[1] - b[1])
  const points: [number, number][] = []
  const n = nums1.length
  for (let i = 0; i < n; i++) {
    points.push([nums1[i], nums2[i]])
  }
  points.sort((a, b) => a[0] - b[0] || a[1] - b[1])
  const allY = new Set<number>()
  points.forEach(([x, y]) => allY.add(y))
  queries.forEach(([x, y]) => allY.add(y))
  const [rank] = sortedSet([...allY])

  // 倒叙遍历点
  // 线段树维护
  const seg = useAtcoderLazySegmentTree(n + 10, {
    e: () => -INF,
    id: () => -INF,
    op: (a, b) => Math.max(a, b),
    mapping: (f, x) => (f === -INF ? x : Math.max(f, x)),
    composition: (f, g) => (f === -INF ? g : Math.max(f, g))
  })

  console.log(points, qWithId)

  let pi = n - 1
  for (let i = qWithId.length - 1; i >= 0; i--) {
    const [qx, qy, qid] = qWithId[i]
    while (pi >= 0 && points[pi][0] >= qx && points[pi][1] >= qy) {
      const [x, y] = points[pi]

      const r = rank(y)
      seg.update(0, r + 1, x + y)
      pi--
    }
    const cur = seg.query(rank(qy), n + 5)
    if (cur !== -INF) res[qid] = cur
  }

  // y值离散化
  return res
}

/**
 * (松)离散化.
 * @returns
 * rank: 给定一个数,返回它的排名`(0-count)`.
 * count: 离散化(去重)后的元素个数.
 */
function sortedSet(nums: number[]): [rank: (num: number) => number, count: number] {
  const allNums = [...new Set(nums)].sort((a, b) => a - b)
  const rank = (num: number) => {
    let left = 0
    let right = allNums.length - 1
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (allNums[mid] >= num) {
        right = mid - 1
      } else {
        left = mid + 1
      }
    }
    return left
  }
  return [rank, allNums.length]
}

// nums1 = [4,3,1,2], nums2 = [2,4,9,5], queries = [[4,1],[1,3],[2,5]]
console.log(
  maximumSumQueries(
    [4, 3, 1, 2],
    [2, 4, 9, 5],
    [
      [4, 1],
      [1, 3],
      [2, 5]
    ]
  )
)
// nums1 = [3,2,5], nums2 = [2,3,4], queries = [[4,4],[3,2],[1,1]]
console.log(
  maximumSumQueries(
    [3, 2, 5],
    [2, 3, 4],
    [
      [4, 4],
      [3, 2],
      [1, 1]
    ]
  )
)
