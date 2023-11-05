/* eslint-disable no-inner-declarations */
/* eslint-disable max-len */
/* eslint-disable class-methods-use-this */

import { UnionFindArrayWithUndoAndWeight } from '../../14_并查集/UnionFindWithUndoAndWeight'

const INF = 2e9 // !超过int32使用2e15

/**
 * 线段树分治undo版.
 * 线段树分治是一种处理动态修改和询问的离线算法.
 * 通过将某一元素的出现时间段在线段树上保存到`log(n)`个结点中,
 * 我们可以 dfs 遍历整棵线段树，运用可撤销数据结构维护来得到每个时间点的答案.
 * @link https://cp-algorithms.com/data_structures/deleting_in_log_n.html
 * @alias OfflineDynamicConnectivity
 */
class SegmentTreeDivideAndConquerUndo {
  private readonly _mutate: (mutationId: number) => void
  private readonly _undo: () => void
  private readonly _query: (queryId: number) => void
  private readonly _mutations: { start: number; end: number; id: number }[] = []
  private readonly _queries: { time: number; id: number }[] = []
  private _nodes: number[][] = [] // 在每个节点上保存对应的变更和查询的编号
  private _mutationId = 0
  private _queryId = 0

  /**
   * dfs 遍历整棵线段树来得到每个时间点的答案.
   * @param mutate 添加编号为`mutationId`的变更后产生的副作用.
   * @param undo 撤销一次`mutate`操作.
   * @param query 响应编号为`queryId`的查询.
   * @complexity 一共调用 **O(nlogn)** 次`mutate`和`undo`，**O(q)** 次`query`.
   */
  constructor(
    mutate: (mutationId: number) => void,
    undo: () => void,
    query: (queryId: number) => void
  )
  constructor(
    options: {
      mutate: (mutationId: number) => void
      undo: () => void
      query: (queryId: number) => void
    } & ThisType<void>
  )
  constructor(arg1: any, arg2?: any, arg3?: any) {
    if (typeof arg1 === 'object') {
      this._mutate = arg1.mutate
      this._undo = arg1.undo
      this._query = arg1.query
    } else {
      this._mutate = arg1
      this._undo = arg2
      this._query = arg3
    }
  }

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

  run(): void {
    if (!this._queries.length) return
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
    this._dfs(1)
  }

  private _dfs(now: number): void {
    const curNodes = this._nodes[now]
    for (let i = 0; i < curNodes.length; i++) {
      const id = curNodes[i]
      if (id & 1) {
        // query
        this._query((id - 1) / 2)
      } else {
        // mutate
        this._mutate(id / 2)
      }
    }

    if (now << 1 < this._nodes.length) this._dfs(now << 1)
    if (((now << 1) | 1) < this._nodes.length) this._dfs((now << 1) | 1)

    // 回溯时撤销
    for (let i = curNodes.length - 1; ~i; i--) {
      if (!(curNodes[i] & 1)) {
        this._undo()
      }
    }
  }

  private _uniqueInplace(sorted: number[]): void {
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

export {
  SegmentTreeDivideAndConquerUndo,
  SegmentTreeDivideAndConquerUndo as OfflineDynamicConnectivity
}

if (require.main === module) {
  testTime()
  function testTime(): void {
    let mutate = 0
    let undo = 0
    let query = 0
    const dc = new SegmentTreeDivideAndConquerUndo({
      mutate(id) {
        mutate++
      },
      undo() {
        undo++
      },
      query(id) {
        query++
      }
    })

    const n = 2e5
    console.time('foo')
    for (let i = 0; i < n; i++) {
      dc.addMutation(0, i)
      dc.addQuery(i)
    }
    dc.run()
    console.timeEnd('foo')
    console.log(mutate, undo, query)
  }

  // Dynamic Graph Vertex Add Component Sum
  // https://judge.yosupo.jp/problem/dynamic_graph_vertex_add_component_sum
  // 0 u v 连接u v (保证u v不连接)
  // 1 u v 断开u v  (保证u v连接)
  // 2 u x 将u的值加上x
  // 3 u 输出u所在连通块的值
  function dynamicGraphVertexAddComponentSum(weights: number[], operations: number[][]): number[] {
    const n = weights.length
    const edges: { u: number; v: number }[] = []
    const existEdge = new Map<number, { id: number; startTime: number }>()
    const adds: { pos: number; add: number }[] = []
    const queries: number[] = []
    const res: number[] = []

    const uf = new UnionFindArrayWithUndoAndWeight(weights, (a, b) => a + b)
    const dc = new SegmentTreeDivideAndConquerUndo({
      mutate(mutationId) {
        if (mutationId >= 0) {
          const e = edges[mutationId]
          uf.union(e.u, e.v)
        } else {
          mutationId = ~mutationId
          const a = adds[mutationId]
          uf.setGroupWeight(a.pos, uf.getGroupWeight(a.pos) + a.add)
        }
      },
      undo() {
        uf.undo()
      },
      query(queryId) {
        const pos = queries[queryId]
        res[queryId] = uf.getGroupWeight(pos)
      }
    })

    for (let time = 0; time < operations.length; time++) {
      const operation = operations[time]
      const op = operation[0]
      if (op === 0) {
        let { 1: u, 2: v } = operation
        if (u < v) {
          u ^= v
          v ^= u
          u ^= v
        }
        const hash = u * n + v
        existEdge.set(hash, { id: edges.length, startTime: time })
        edges.push({ u, v })
      } else if (op === 1) {
        let { 1: u, 2: v } = operation
        if (u < v) {
          u ^= v
          v ^= u
          u ^= v
        }
        const hash = u * n + v
        const item = existEdge.get(hash)!
        dc.addMutation(item.startTime, time, item.id)
        existEdge.delete(hash)
      } else if (op === 2) {
        const { 1: pos, 2: add } = operation
        const id = ~adds.length
        dc.addMutation(time, INF, id)
        adds.push({ pos, add })
      } else {
        const pos = operation[1]
        dc.addQuery(time, queries.length)
        queries.push(pos)
        res.push(0)
      }
    }

    existEdge.forEach(item => {
      dc.addMutation(item.startTime, INF, item.id)
    })

    dc.run()

    return res
  }

  const weights = [1, 10, 100, 1000, 10000]
  const operations = [
    [0, 0, 1],
    [0, 1, 2],
    [0, 2, 3],
    [0, 3, 4],
    [0, 0, 4],
    [3, 3],
    [1, 1, 2],
    [3, 1],
    [1, 3, 4],
    [3, 0],
    [2, 1, 100000],
    [3, 1],
    [0, 1, 4],
    [3, 2],
    [0, 3, 4],
    [3, 0]
  ]

  // 11111
  // 11111
  // 10011
  // 110011
  // 1100
  // 111111
  console.log(dynamicGraphVertexAddComponentSum(weights, operations))
}
