/* eslint-disable no-inner-declarations */
/* eslint-disable space-in-parens */
/* eslint-disable prefer-destructuring */
/* eslint-disable max-len */

type Edge = {
  from: number
  to: number
  capacity: number
  flow: number
  cost: number
  id: number
}

type Neighbor = {
  to: number
  /** 反向边在邻接表中的下标. */
  rid: number
  /** 边的残量. */
  capacity: number
  cost: number
  /** -1表示是反向边. */
  id: number
}

const INF = 2e15

class MinCostMaxFlow {
  private readonly _graph: Neighbor[][]
  private readonly _n: number
  private readonly _start: number
  private readonly _end: number
  private readonly _dist: number[]
  private readonly _pre: { v: number; id: number }[]
  private _ei = 0

  constructor(n: number, start: number, end: number) {
    this._n = n
    this._start = start
    this._end = end
    this._graph = Array(n)
    this._dist = Array(n)
    this._pre = Array(n)
    for (let i = 0; i < n; i++) {
      this._graph[i] = []
      this._pre[i] = {} as any
    }
  }

  addEdge(from: number, to: number, capacity: number, cost: number): void {
    this._graph[from].push({ to, rid: this._graph[to].length, capacity, cost, id: this._ei })
    this._graph[to].push({
      to: from,
      rid: this._graph[from].length - 1,
      capacity: 0,
      cost: -cost,
      id: -1
    })
    this._ei++
  }

  flow(limit = INF): [maxFlow: number, minCost: number] {
    let maxFlow = 0
    let minCost = 0
    while (maxFlow < limit) {
      if (!this._spfa()) break
      let flow = INF
      for (let cur = this._end; cur !== this._start; ) {
        const { v, id } = this._pre[cur]
        const edge = this._graph[v][id]
        if (edge.capacity < flow) flow = edge.capacity
        cur = v
      }
      for (let cur = this._end; cur !== this._start; ) {
        const { v, id } = this._pre[cur]
        const edge = this._graph[v][id]
        edge.capacity -= flow
        this._graph[cur][edge.rid].capacity += flow
        cur = v
      }
      maxFlow += flow
      minCost += this._dist[this._end] * flow
    }
    return [maxFlow, minCost]
  }

  /**
   * @returns (flow, cost) 的每个转折点.
   */
  slope(limit = INF): [flow: number, cost: number][] {
    const res: [flow: number, cost: number][] = []
    let maxFlow = 0
    let minCost = 0
    while (maxFlow < limit) {
      if (!this._spfa()) break
      let flow = INF
      for (let cur = this._end; cur !== this._start; ) {
        const { v, id } = this._pre[cur]
        const edge = this._graph[v][id]
        if (edge.capacity < flow) flow = edge.capacity
        cur = v
      }
      for (let cur = this._end; cur !== this._start; ) {
        const { v, id } = this._pre[cur]
        const edge = this._graph[v][id]
        edge.capacity -= flow
        this._graph[cur][edge.rid].capacity += flow
        cur = v
      }
      maxFlow += flow
      minCost += this._dist[this._end] * flow
      res.push([maxFlow, minCost])
    }
    return res
  }

  /**
   * @warning 注意根据from,to排除虚拟源点汇点; `flow>0` 才是流经的边.
   */
  getEdges(): Edge[] {
    const res: Edge[] = []
    for (let from = 0; from < this._n; from++) {
      const nexts = this._graph[from]
      for (let i = 0; i < nexts.length; i++) {
        const { to, capacity, cost, id, rid } = nexts[i]
        if (id === -1) continue
        const tos = this._graph[to]
        res.push({
          from,
          to,
          cost,
          id,
          capacity: capacity + tos[rid].capacity,
          flow: tos[rid].capacity
        })
      }
    }
    return res
  }

  private _spfa(): boolean {
    const { _start, _end, _dist, _pre, _n, _graph } = this
    _dist.fill(INF)
    _dist[_start] = 0
    const inQueue = new Uint8Array(_n)
    inQueue[_start] = 1
    let queue = [_start]
    while (queue.length) {
      const nextQueue: number[] = []
      for (let i = 0; i < queue.length; i++) {
        const cur = queue[i]
        inQueue[cur] = 0
        const nexts = _graph[cur]
        for (let j = 0; j < nexts.length; j++) {
          const { to, capacity, cost } = nexts[j]
          if (!capacity) continue
          if (_dist[cur] + cost < _dist[to]) {
            _dist[to] = _dist[cur] + cost
            _pre[to] = { v: cur, id: j }
            if (!inQueue[to]) {
              nextQueue.push(to)
              inQueue[to] = 1
            }
          }
        }
      }
      queue = nextQueue
    }
    return _dist[_end] < INF
  }
}

export { MinCostMaxFlow }

if (require.main === module) {
  // https://leetcode.cn/problems/maximum-and-sum-of-array/
  function maximumANDSum(nums: number[], numSlots: number): number {
    const n = nums.length
    const m = numSlots
    const [START, END] = [n + m + 2, n + m + 3]
    const mcmf = new MinCostMaxFlow(n + m + 5, START, END)
    for (let i = 0; i < n; i++) {
      for (let j = 0; j < numSlots; j++) {
        mcmf.addEdge(i, j + n, 1, -(nums[i] & (j + 1)))
      }
    }

    for (let i = 0; i < n; i++) mcmf.addEdge(START, i, 1, 0)
    for (let i = 0; i < numSlots; i++) mcmf.addEdge(i + n, END, 2, 0)
    return -mcmf.flow()[1]
  }

  // https://leetcode.cn/problems/minimum-moves-to-spread-stones-over-grid/solutions/2435374/zui-xiao-fei-yong-zui-da-liu-by-98137766-afei/
  // 将石头分散到网格图的最少移动次数
  // 1.源点向所有格子连一条容量为 grid[r][c]grid[r][c]grid[r][c], 费用为000的边；
  // 2.每个格子的向其他格子连一条容量为INFINFINF, 费用为曼哈顿距离曼哈顿距离曼哈顿距离的边；
  // 3.每个格子向汇点连一条容量为111, 费用为000的边；
  // 4.取得最大流时的最小费用即为答案。
  function minimumMoves(grid: number[][]): number {
    const n = grid.length
    const OFFSET = n * n
    const START = 2 * OFFSET
    const END = START + 1
    const mcmf = new MinCostMaxFlow(END + 1, START, END)
    const hash = (r: number, c: number) => r * n + c

    for (let r = 0; r < n; r++) {
      for (let c = 0; c < n; c++) {
        mcmf.addEdge(START, hash(r, c), grid[r][c], 0)
      }
    }

    for (let r1 = 0; r1 < n; r1++) {
      for (let c1 = 0; c1 < n; c1++) {
        const cur = hash(r1, c1)
        for (let r2 = 0; r2 < n; r2++) {
          for (let c2 = 0; c2 < n; c2++) {
            const next = hash(r2, c2)
            mcmf.addEdge(cur, OFFSET + next, INF, Math.abs(r1 - r2) + Math.abs(c1 - c2))
          }
        }
      }
    }

    for (let r = 0; r < n; r++) {
      for (let c = 0; c < n; c++) {
        mcmf.addEdge(OFFSET + hash(r, c), END, 1, 0)
      }
    }

    return mcmf.flow()[1]
  }
}
