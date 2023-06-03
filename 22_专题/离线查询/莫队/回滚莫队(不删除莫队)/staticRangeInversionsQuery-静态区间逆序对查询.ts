import { FenwickTree } from '../../../../6_tree/树状数组/经典题/BIT-ei1333'
import { MoRollback } from './MoRollback'

/**
 * 静态区间逆序对查询.
 */
function staticRangeInversionsQuery(
  arr: number[],
  queries: [start: number, end: number][]
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
  const history: number[] = []
  const bit = new FenwickTree(rank.size) // 值域树状数组

  mo.run(add, reset, snapshot, rollback, query)
  return res

  // TODO
  function add(index: number, delta: -1 | 1): void {
    if (delta === 1) {
    } else {
    }
  }
  function _move(state: number): void {
    while (history.length > state) {
      const x = history.pop()!
      bit.add(x, -1)
    }
  }

  function reset(): void {
    _move(0)
    cur = 0
  }
  function snapshot(): void {
    snap = history.length
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
