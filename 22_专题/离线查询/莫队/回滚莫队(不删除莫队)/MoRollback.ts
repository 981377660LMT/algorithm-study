/* eslint-disable no-inner-declarations */

/**
 * 回滚莫队(不删除莫队),复杂度和普通莫队一样.
 * !删除操作很麻烦的时候使用.
 */
class MoRollback {
  private readonly _chunkSize: number
  private readonly _left: Uint32Array
  private readonly _right: Uint32Array
  private readonly _order: Uint32Array
  private _leftPtr = 0
  private _rightPtr = 0

  constructor(n: number, q: number) {
    const chunkSize = Math.max(1, (n / Math.max(1, Math.sqrt((q * 2) / 3) | 0)) | 0)
    const order = new Uint32Array(q)
    for (let i = 0; i < q; ++i) order[i] = i
    this._chunkSize = chunkSize
    this._left = new Uint32Array(q)
    this._right = new Uint32Array(q)
    this._order = order
  }

  /**
   * 添加一个查询，查询范围为左闭右开区间`[start, end)`.
   * 0 <= start <= end <= n.
   */
  addQuery(start: number, end: number): void {
    this._left[this._leftPtr++] = start
    this._right[this._rightPtr++] = end
  }

  /**
   * 返回每个查询的结果.
   * @param add 将数据添加到窗口.
   * @param reset 将窗口重置为初始状态.
   * @param snapshot 保存当前窗口的状态.
   * @param rollback 恢复窗口的状态.
   * @param query 查询窗口的状态.
   */
  run(
    add: (index: number, delta: -1 | 1) => void,
    reset: () => void,
    snapshot: () => void,
    rollback: () => void,
    query: (qi: number) => void
  ): void {
    const left = this._left
    const right = this._right
    const order = this._order
    const chunkSize = this._chunkSize
    order.sort((o1, o2) => {
      const iblock = (left[o1] / chunkSize) | 0
      const jblock = (left[o2] / chunkSize) | 0
      return iblock - jblock || right[o1] - right[o2]
    })

    reset()
    for (let i = 0; i < order.length; ++i) {
      const index = order[i]
      if (right[index] - left[index] < chunkSize) {
        for (let j = left[index]; j < right[index]; ++j) add(j, 1)
        query(index)
        rollback()
      }
    }

    let nr = 0
    let lastBlock = -1
    for (let i = 0; i < order.length; ++i) {
      const index = order[i]
      if (right[index] - left[index] < chunkSize) continue
      const block = (left[index] / chunkSize) | 0
      if (lastBlock !== block) {
        reset()
        lastBlock = block
        nr = (block + 1) * chunkSize
      }
      while (nr < right[index]) add(nr++, 1)
      snapshot()
      for (let j = (block + 1) * chunkSize - 1; j >= left[index]; --j) add(j, -1)
      query(index)
      rollback()
    }
  }
}

export { MoRollback }

if (require.main === module) {
  // 历史的研究
  // 给定一个数组nums和q个查询(l,r)
  // 每次查询[l,r]区间内的`重要度`,一个数字num的重要度定义为`num乘以区间内num的个数`
  // https://www.luogu.com.cn/problem/AT_joisc2014_c
  // https://atcoder.jp/contests/joisc2014/tasks/joisc2014_c
  function solve1(nums: number[], queries: [start: number, end: number][]): number[] {
    const n = nums.length
    const q = queries.length
    const mo = new MoRollback(n, q)
    queries.forEach(([start, end]) => mo.addQuery(start, end))

    const res: number[] = Array(q)
    let cur = 0 // 当前区间的答案
    let snap = 0 // 当前区间的快照
    let snapCur = 0 // 当前快照的答案
    const history: number[] = []
    const counter = new Map<number, number>()
    mo.run(add, reset, snapshot, rollback, query)
    return res

    // TODO
    function add(index: number, _delta: -1 | 1): void {
      const x = nums[index]
      const count = counter.get(x) || 0
      counter.set(x, count + 1)
      cur = Math.max(cur, x * (count + 1))
      history.push(x)
    }
    function _move(state: number): void {
      while (history.length > state) {
        const x = history.pop()!
        counter.set(x, (counter.get(x) || 0) - 1) // TODO
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

  // 区间内相同的数的最远距离
  // 给定一个序列，多次询问一段区间 [l,r]，求区间中相同的数的最远间隔距离。
  // 如果区间内不存在两个数相同，则输出 0。
  // 序列中两个元素的间隔距离指的是两个元素下标差的绝对值。
  //
  // !维护每个数在区间内索引的最大值和最小值.
  // https://www.luogu.com.cn/problem/P5906
  function solve2(nums: number[], queries: [number, number][]): number[] {
    const n = nums.length
    const q = queries.length
    const mo = new MoRollback(n, q)
    queries.forEach(([start, end]) => mo.addQuery(start, end))

    // 哈希
    const _pool = new Map<unknown, number>()
    const id = (o: unknown): number => {
      if (!_pool.has(o)) _pool.set(o, _pool.size)
      return _pool.get(o)!
    }
    for (let i = 0; i < n; ++i) {
      nums[i] = id(nums[i])
    }

    const res: number[] = Array(q)
    let cur = 0 // 当前区间的答案
    let snap = 0 // 当前区间的快照
    let snapCur = 0 // 当前快照的答案
    const history: [x: number, preMinPos: number, preMaxPos: number][] = []
    const minPos = new Int32Array(n)
    const maxPos = new Int32Array(n)
    for (let i = 0; i < n; ++i) {
      minPos[i] = n
      maxPos[i] = -1
    }
    mo.run(add, reset, snapshot, rollback, query)
    return res

    // TODO
    function add(index: number, _delta: -1 | 1): void {
      const x = nums[index]
      const preMinPos = minPos[x]
      const preMaxPos = maxPos[x]
      minPos[x] = Math.min(minPos[x], index)
      maxPos[x] = Math.max(maxPos[x], index)
      cur = Math.max(cur, maxPos[x] - minPos[x])
      history.push([x, preMinPos, preMaxPos])
    }
    function _move(state: number): void {
      while (history.length > state) {
        const [x, preMinPos, preMaxPos] = history.pop()!
        minPos[x] = preMinPos
        maxPos[x] = preMaxPos
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
      if (cur > 0) res[qi] = cur
    }
  }
}
