import assert from 'assert'

import { FenwickTree } from '../../../../6_tree/树状数组/经典题/BIT-ei1333'
import { MoRollback } from './MoRollback'

/**
 * 静态区间逆序对查询.
 * 时间复杂度：O(n√nlogn).
 * n,q 1e5 => 2s.
 */
function staticRangeInversionsQuery(
  arr: number[],
  queries: [start: number, end: number][] | number[][]
): number[] {
  const n = arr.length
  const q = queries.length
  const mo = new MoRollback(n, q)
  queries.forEach(([start, end]) => mo.addQuery(start, end))

  // 离散化
  const allNums = [...new Set(arr)].sort((a, b) => a - b)
  const rank = new Map<number, number>()
  for (let i = 0; i < allNums.length; ++i) rank.set(allNums[i], i)
  const newNums = new Uint32Array(n)
  for (let i = 0; i < n; ++i) newNums[i] = rank.get(arr[i])!

  let cur = 0 // 当前区间的答案
  let snap = 0 // 当前区间的快照
  let snapCur = 0 // 当前快照的答案
  const res: number[] = Array(q)
  const history = new Uint32Array(n)
  let historyPtr = 0
  const bit = new FenwickTree(rank.size) // 值域树状数组

  mo.run(add, reset, snapshot, rollback, query)
  return res

  // TODO
  function add(index: number, delta: -1 | 1): void {
    if (delta === 1) {
      const x = newNums[index]
      cur += bit.queryRange(x + 1, rank.size)
      bit.add(x, 1)
      history[historyPtr++] = x
    } else {
      const x = newNums[index]
      cur += bit.query(x)
      bit.add(x, 1)
      history[historyPtr++] = x
    }
  }
  function _move(state: number): void {
    while (historyPtr > state) {
      const x = history[--historyPtr]
      bit.add(x, -1)
    }
  }

  function reset(): void {
    _move(0)
    cur = 0
  }
  function snapshot(): void {
    snap = historyPtr
    snapCur = cur
  }
  function rollback(): void {
    _move(snap)
    cur = snapCur
  }
  function query(qi: number): void {
    res[qi] = cur
  }
}

export {}
if (require.main === module) {
  //   4 2
  // 4 1 4 0
  // 1 3
  // 0 4
  const nums = [4, 1, 4, 0]
  const queries = [
    [1, 3],
    [0, 4]
  ]
  assert.deepEqual(staticRangeInversionsQuery(nums, queries), [0, 4])

  // perf
  const n = 1e5
  const q = 1e5
  const nums2 = Array(n)
    .fill(0)
    .map((_, i) => i)
  const queries2 = Array(q)
  for (let i = 0; i < q; ++i) {
    const start = Math.floor(Math.random() * n)
    const end = Math.floor(Math.random() * n)
    queries2[i] = [Math.min(start, end), Math.max(start, end)]
  }
  console.time('perf')
  staticRangeInversionsQuery(nums2, queries2)
  console.timeEnd('perf') // 2.051s
}
