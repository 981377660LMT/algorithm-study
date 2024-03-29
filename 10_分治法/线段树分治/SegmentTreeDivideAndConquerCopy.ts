/* eslint-disable max-len */
/* eslint-disable no-inner-declarations */
/* eslint-disable class-methods-use-this */

/**
 * 线段树分治copy版.
 * 如果修改操作难以撤销，可以在每个节点处保存一份副本.
 * !调用O(n)次拷贝注意不要超出内存.
 */
class SegmentTreeDivideAndConquerCopy {
  private _initState!: unknown
  private _mutate!: (state: unknown, mutationId: number) => void
  private _copy!: (state: unknown) => unknown
  private _query!: (state: unknown, queryId: number) => void
  private _undo?: () => void
  private readonly _mutations: { start: number; end: number; id: number }[] = []
  private readonly _queries: { time: number; id: number }[] = []
  private _nodes: number[][] = [] // 在每个节点上保存对应的变更和查询的编号
  private _mutationId = 0
  private _queryId = 0

  /**
   * 在时间范围`[startTime, endTime)`内添加一个编号为`id`的变更.
   */
  addMutation(startTime: number, endTime: number, id?: number): void {
    if (startTime >= endTime) return
    if (id == undefined) id = this._mutationId++
    this._mutations.push({ start: startTime, end: endTime, id })
  }

  /**
   * 在时间`time`时添加一个编号为`id`的查询.
   */
  addQuery(time: number, id?: number): void {
    if (id == undefined) id = this._queryId++
    this._queries.push({ time, id })
  }

  /**
   * dfs 遍历整棵线段树来得到每个时间点的答案.
   * @param initState 数据结构的初始状态.
   * @param mutate 添加编号为`mutationId`的变更后产生的副作用.
   * @param copy 拷贝一份数据结构的副本.
   * @param query 响应编号为`queryId`的查询.
   * @param undo 撤销最近一次变更.
   * @complexity 一共调用 **O(nlogn)** 次`mutate`，**O(n)** 次`copy` , **O(q)** 次`query`，**O(nlogn)** 次`undo`.
   */
  run<S>(
    initState: S,
    options: {
      mutate: (state: S, mutationId: number) => void
      copy: (state: S) => S
      query: (state: S, queryId: number) => void
      undo?: () => void
    } & ThisType<void>
  ): void
  run<S>(
    initState: S,
    mutate: (state: S, mutationId: number) => void,
    copy: (state: S) => S,
    query: (state: S, queryId: number) => void,
    undo?: () => void
  ): void
  run(arg1: any, arg2: any, arg3?: any, arg4?: any, arg5?: any): void {
    if (!this._queries.length) return
    this._initState = arg1
    if (typeof arg2 === 'object') {
      this._mutate = arg2.mutate
      this._copy = arg2.copy
      this._query = arg2.query
      this._undo = arg2.undo
    } else {
      this._mutate = arg2
      this._copy = arg3
      this._query = arg4
      this._undo = arg5
    }
    const times: number[] = Array(this._queries.length)
    for (let i = 0; i < this._queries.length; i++) times[i] = this._queries[i].time
    times.sort((a, b) => a - b)
    this._uniqueInplace(times)
    const usedTimes = new Uint8Array(times.length + 1)
    usedTimes[0] = 1
    for (let i = 0; i < this._mutations.length; i++) {
      const e = this._mutations[i]
      usedTimes[this._lowerBound(times, e.start)] = 1
      usedTimes[this._lowerBound(times, e.end)] = 1
    }
    for (let i = 1; i < times.length; i++) {
      if (!usedTimes[i]) times[i] = times[i - 1]
    }
    this._uniqueInplace(times)

    const n = times.length
    let offset = 1
    while (offset < n) offset <<= 1
    this._nodes = Array(offset + n)
    for (let i = 0; i < this._nodes.length; i++) this._nodes[i] = []
    for (let i = 0; i < this._mutations.length; i++) {
      const e = this._mutations[i]
      let left = offset + this._lowerBound(times, e.start)
      let right = offset + this._lowerBound(times, e.end)
      const eid = e.id * 2
      while (left < right) {
        if (left & 1) this._nodes[left++].push(eid) // mutate
        if (right & 1) this._nodes[--right].push(eid)
        left >>>= 1
        right >>>= 1
      }
    }

    for (let i = 0; i < this._queries.length; i++) {
      const q = this._queries[i]
      this._nodes[offset + this._upperBound(times, q.time) - 1].push(q.id * 2 + 1) // query
    }
    this._dfs(1, this._initState)
  }

  private _dfs(now: number, state: unknown): void {
    const curNodes = this._nodes[now]
    for (let i = 0; i < curNodes.length; i++) {
      const id = curNodes[i]
      if (id & 1) {
        this._query(state, (id - 1) / 2)
      } else {
        this._mutate(state, id / 2)
      }
    }

    if (now << 1 < this._nodes.length) {
      this._dfs(now << 1, this._copy(state))
    }
    if (((now << 1) | 1) < this._nodes.length) {
      this._dfs((now << 1) | 1, this._copy(state))
    }

    if (this._undo) {
      for (let i = 0; i < curNodes.length; i++) {
        const id = curNodes[i]
        if (!(id & 1)) {
          this._undo()
        }
      }
    }
  }

  private _uniqueInplace(sorted: number[]): void {
    if (sorted.length === 0) return
    let slow = 0
    for (let fast = 0; fast < sorted.length; fast++) {
      if (sorted[fast] !== sorted[slow]) sorted[++slow] = sorted[fast]
    }
    sorted.length = slow + 1
  }

  private _lowerBound(arr: ArrayLike<number>, target: number): number {
    let left = 0
    let right = arr.length - 1
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (arr[mid] < target) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }

  private _upperBound(arr: ArrayLike<number>, target: number): number {
    let left = 0
    let right = arr.length - 1
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (arr[mid] <= target) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }
}

export { SegmentTreeDivideAndConquerCopy }

if (require.main === module) {
  // 238. 除自身以外数组的乘积
  // https://leetcode.cn/problems/product-of-array-except-self/
  function productExceptSelf(nums: number[]): number[] {
    const n = nums.length
    const res = Array(n).fill(1)
    const seg = new SegmentTreeDivideAndConquerCopy()

    // 第i次变更在时间段 `[0, i) + [i+1, n)` 内存在.
    for (let i = 0; i < n; i++) {
      seg.addMutation(0, i, i)
      seg.addMutation(i + 1, n, i)
    }
    for (let i = 0; i < n; i++) {
      seg.addQuery(i, i)
    }
    seg.run(
      { value: 1 },
      {
        mutate(state, mutationId) {
          state.value *= nums[mutationId]
        },
        copy(state) {
          return { value: state.value }
        },
        query(state, queryId) {
          res[queryId] = state.value
        }
      }
    )

    return res
  }
}
