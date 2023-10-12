/* eslint-disable max-len */

// https://www.luogu.com.cn/problem/P6578
// P6578 [Ynoi2019]魔法少女网站
//

// 给定一个数组nums和一些查询(start,end,x)，对每个查询回答区间[start,end)内有多少个子数组最大值不超过x.
// nums.length<=1e5,nums[i]<=1e5,查询个数<=1e5
// 离线查询+线段树.
// 将所有查询按照x从小到大排序,然后从小到大依次处理，即维护一个01数组.
// 长为len的极长连续段的贡献为len*(len+1)/2.
// 用线段树维护.

import { createPointSetRangeLongestOne } from '../template/atcoder_segtree/SegmentTreeUtils'

function 魔法少女网站无修改版(nums: number[], queries: [start: number, end: number, threshold: number][]): number[] {
  type Query = { start: number; end: number; threshold: number; id: number }
  type Pair = { value: number; index: number }

  const n = nums.length
  const q = queries.length

  const sortedQueries: Query[] = Array(q)
  for (let i = 0; i < q; i++) {
    const { 0: start, 1: end, 2: threshold } = queries[i]
    sortedQueries[i] = { start, end, threshold, id: i }
  }
  sortedQueries.sort((a, b) => a.threshold - b.threshold)

  const sortedPairs: Pair[] = Array(n)
  for (let i = 0; i < n; i++) sortedPairs[i] = { value: nums[i], index: i }
  sortedPairs.sort((a, b) => a.value - b.value)

  const res: number[] = Array(q).fill(0)
  const { tree, fromElement } = createPointSetRangeLongestOne(n)
  let ptr = 0
  sortedQueries.forEach(({ start, end, threshold, id }) => {
    while (ptr < n && sortedPairs[ptr].value <= threshold) {
      tree.set(sortedPairs[ptr].index, fromElement(1))
      ptr++
    }
    res[id] = tree.query(start, end).pairCount
  })
  return res
}

if (require.main === module) {
  const nums = [1, 2, 3, 4, 5]
  const queries = [
    [1, 5, 3],
    [1, 5, 4],
    [1, 5, 5],
    [1, 5, 6],
    [1, 5, 7]
  ]
  console.info(魔法少女网站无修改版(nums, queries))
}
