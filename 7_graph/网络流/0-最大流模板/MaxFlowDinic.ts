/* eslint-disable no-inner-declarations */
/* eslint-disable max-len */

type Edge = {
  to: number
  /** 剩余容量. */
  capicity: number
  /** 反向边在 G[to] 中的序号. */
  rev: number
  isRev: boolean
  index: number
}

const INF = 2e15

/**
 * Dinic算法求最大流.
 * @complexity `O(V^2*E)` , 二分图上为 O(E*√V)
 */
class MaxFlowDinic {
  private readonly _graph: Edge[][]
  private readonly _minCost: Int32Array
  private readonly _iter: Uint32Array
  private readonly _n: number
  private readonly _visitedEdge: Set<number> = new Set()

  constructor(n: number) {
    this._n = n
    this._graph = Array(n)
    this._minCost = new Int32Array(n)
    this._iter = new Uint32Array(n)
    for (let i = 0; i < n; i++) this._graph[i] = []
  }

  /**
   * 内部会对边去重.
   */
  addEdge(from: number, to: number, cap: number, index = -1): void {
    const hash = from * this._n + to
    if (this._visitedEdge.has(hash)) return
    this._visitedEdge.add(hash)
    this._graph[from].push({ to, capicity: cap, rev: this._graph[to].length, isRev: false, index })
    this._graph[to].push({ to: from, capicity: 0, rev: this._graph[from].length - 1, isRev: true, index })
  }

  maxFlow(start: number, target: number): number {
    let flow = 0
    while (this._buildAugmentingPath(start, target)) {
      this._iter.fill(0)
      let f = 0
      while (true) {
        f = this._findMinDistAugmentPath(start, target, INF)
        if (!f) break
        flow += f
      }
    }
    return flow
  }

  /**
   * @returns (from,to,流量,容量).
   * flow = revEdge.cap; cap = e.cap + revEdge.cap.
   */
  getEdges(): { from: number; to: number; flow: number; cap: number }[] {
    const res: { from: number; to: number; flow: number; cap: number }[] = []
    for (let i = 0; i < this._n; i++) {
      this._graph[i].forEach(e => {
        if (e.isRev) return
        const revEdge = this._graph[e.to][e.rev]
        res.push({ from: i, to: e.to, flow: revEdge.capicity, cap: e.capicity + revEdge.capicity })
      })
    }
    return res
  }

  private _findMinDistAugmentPath(index: number, target: number, flow: number): number {
    if (index === target) return flow
    const nexts = this._graph[index]
    let ptr = this._iter[index]
    while (ptr < nexts.length) {
      const e = nexts[ptr]
      if (e.capicity > 0 && this._minCost[index] < this._minCost[e.to]) {
        const f = this._findMinDistAugmentPath(e.to, target, Math.min(flow, e.capicity))
        if (f > 0) {
          nexts[ptr].capicity -= f
          this._graph[e.to][e.rev].capicity += f
          return f
        }
      }
      ptr++
      this._iter[index]++
    }
    return 0
  }

  private _buildAugmentingPath(start: number, target: number): boolean {
    this._minCost.fill(-1)
    this._minCost[start] = 0
    let queue = [start]
    while (queue.length) {
      const nextQueue: number[] = []
      for (let i = 0; i < queue.length; i++) {
        const cur = queue[i]
        const nexts = this._graph[cur]
        for (let j = 0; j < nexts.length; j++) {
          const edge = nexts[j]
          if (edge.capicity > 0 && this._minCost[edge.to] === -1) {
            this._minCost[edge.to] = this._minCost[cur] + 1
            nextQueue.push(edge.to)
          }
        }
      }
      queue = nextQueue
    }

    return this._minCost[target] !== -1
  }
}

export { MaxFlowDinic }

if (require.main === module) {
  function maximumInvitations(grid: number[][]): number {
    const [ROW, COL] = [grid.length, grid[0].length]
    const START = ROW + COL
    const END = START + 1
    const F = new MaxFlowDinic(ROW + COL + 2)
    for (let r = 0; r < ROW; r++) {
      for (let c = 0; c < COL; c++) {
        if (grid[r][c] === 1) {
          F.addEdge(START, r, 1)
          F.addEdge(r, ROW + c, 1)
          F.addEdge(ROW + c, END, 1)
        }
      }
    }
    return F.maxFlow(START, END)
  }

  console.log(
    maximumInvitations([
      [1, 1, 1],
      [1, 0, 1],
      [0, 0, 1]
    ])
  )
}
