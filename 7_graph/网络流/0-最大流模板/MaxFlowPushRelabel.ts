/* eslint-disable no-inner-declarations */
/* eslint-disable no-labels */
/* eslint-disable no-unreachable-loop */

import { HeapUint32 } from '../../../8_heap/SiftHeap'

type Edge = {
  to: number
  cap: number
  rid: number
}

const INF = 2e15

class MaxFlowPushRelabel {
  private readonly _graph: Edge[][]
  private readonly _n: number
  private readonly _visitedEdge: Set<number> = new Set()

  constructor(n: number) {
    this._n = n
    this._graph = Array(n)
    for (let i = 0; i < n; i++) this._graph[i] = []
  }

  addEdge(from: number, to: number, cap: number): void {
    const hash = from * this._n + to
    if (this._visitedEdge.has(hash)) return
    this._visitedEdge.add(hash)
    this._graph[from].push({ to, cap, rid: this._graph[to].length })
    this._graph[to].push({ to: from, cap: 0, rid: this._graph[from].length - 1 })
  }

  maxFlow(start: number, target: number): number {
    const n = this._graph.length
    const dist = new Int32Array(n).fill(-1)
    dist[target] = 0
    const distCounter = new Uint32Array(2 * n)
    let queue: number[] = [target]
    while (queue.length) {
      const nextQueue: number[] = []
      for (let i = 0; i < queue.length; i++) {
        const v = queue[i]
        distCounter[dist[v]]++
        const nexts = this._graph[v]
        for (let j = 0; j < nexts.length; j++) {
          const e = nexts[j]
          if (dist[e.to] < 0) {
            dist[e.to] = dist[v] + 1
            nextQueue.push(e.to)
          }
        }
      }
      queue = nextQueue
    }
    dist[start] = n

    const exFlow = Array(n).fill(0)
    const pq = new HeapUint32(n, (i, j) => dist[i] > dist[j])
    const inQueue = new Uint8Array(n)

    const push = (v: number, f: number, e: Edge): void => {
      const w = e.to
      e.cap -= f
      this._graph[w][e.rid].cap += f
      exFlow[v] -= f
      exFlow[w] += f
      if (w !== start && w !== target && !inQueue[w]) {
        pq.push(w)
        inQueue[w] = 1
      }
    }

    for (let i = 0; i < this._graph[start].length; i++) {
      const e = this._graph[start][i]
      if (e.cap > 0) push(start, e.cap, e)
    }

    while (pq.size) {
      const v = pq.pop()!
      inQueue[v] = 0
      outer: while (true) {
        const nexts = this._graph[v]
        for (let i = 0; i < nexts.length; i++) {
          const e = nexts[i]
          if (e.cap > 0 && dist[e.to] < dist[v]) {
            push(v, Math.min(e.cap, exFlow[v]), e)
            if (!exFlow[v]) break outer
          }
        }

        const dv = dist[v]
        if (dv !== -1) {
          if (--distCounter[dv] === 0) {
            for (let i = 0; i < n; i++) {
              if (i !== start && i !== target && dv < dist[i] && dist[i] <= n) {
                dist[i] = n + 1
              }
            }
          }
        }

        let minD = INF
        for (let i = 0; i < nexts.length; i++) {
          const e = nexts[i]
          if (e.cap > 0 && dist[e.to] < minD) {
            minD = dist[e.to]
          }
        }
        dist[v] = minD + 1
        distCounter[dist[v]]++
      }
    }

    return exFlow[target]
  }
}

export { MaxFlowPushRelabel }

if (require.main === module) {
  function maximumInvitations(grid: number[][]): number {
    const [ROW, COL] = [grid.length, grid[0].length]
    const START = ROW + COL
    const END = START + 1
    const F = new MaxFlowPushRelabel(ROW + COL + 2)
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
