/* eslint-disable no-inner-declarations */
/* eslint-disable max-len */

type Edge = {
  to: number
  cap: number
  rev: number
  isRev: boolean
  index: number
}

/**
 * 最高标号预流推进算法.
 */
class MaxFlowPushRelabel {
  /** 如果总的流量与边的流量小于 2^32, 内部使用 Uint32Array. */
  private _exFlow!: number[] | Uint32Array
  private readonly _graph: Edge[][]
  private readonly _potential: Uint32Array
  private readonly _curEdge: Uint32Array
  private readonly _allVertex: _List
  private readonly _activeVertex: _Stack
  private readonly _n: number
  private readonly _visitedEdge: Set<number> = new Set()
  private _height = -1
  private _relabels = 0
  private _capSum = 0

  constructor(n: number) {
    this._n = n
    this._graph = Array(n)
    for (let i = 0; i < n; i++) this._graph[i] = []
    this._potential = new Uint32Array(n)
    this._curEdge = new Uint32Array(n)
    this._allVertex = new _List(n, n)
    this._activeVertex = new _Stack(n, n)
  }

  /**
   * 内部会对边去重.
   */
  addEdge(from: number, to: number, cap: number, index = -1): void {
    const hash = from * this._n + to
    if (this._visitedEdge.has(hash)) return
    this._visitedEdge.add(hash)
    this._graph[from].push({ to, cap, rev: this._graph[to].length, isRev: false, index })
    this._graph[to].push({ to: from, cap: 0, rev: this._graph[from].length - 1, isRev: true, index })
    this._capSum += cap
  }

  maxFlow(start: number, target: number): number {
    let level = this._init(start, target)
    while (level >= 0) {
      if (this._activeVertex.empty(level)) {
        level--
        continue
      }
      const u = this._activeVertex.top(level)
      this._activeVertex.pop(level)
      level = this._discharge(u, target)
      if (this._relabels * 2 >= this._n) {
        level = this._globalRelabel(target)
        this._relabels = 0
      }
    }
    return this._exFlow[target]
  }

  private _calcActive(t: number): number {
    const n = this._n
    this._height = -1
    for (let i = 0; i < n; i++) {
      if (this._potential[i] < n) {
        this._curEdge[i] = 0
        this._height = Math.max(this._height, this._potential[i])
        this._allVertex.insert(this._potential[i], i)
        if (this._exFlow[i] > 0 && i !== t) this._activeVertex.push(this._potential[i], i)
      } else {
        this._potential[i] = n + 1
      }
    }
    return this._height
  }

  private _bfs(t: number): void {
    const n = this._n
    for (let i = 0; i < n; i++) {
      this._potential[i] = Math.max(this._potential[i], n)
    }
    this._potential[t] = 0
    let queue: number[] = [t]
    while (queue.length) {
      const nextQueue: number[] = []
      for (let i = 0; i < queue.length; i++) {
        const p = queue[i]
        const nexts = this._graph[p]
        for (let j = 0; j < nexts.length; j++) {
          const e = nexts[j]
          if (this._potential[e.to] === n && this._graph[e.to][e.rev].cap > 0) {
            this._potential[e.to] = this._potential[p] + 1
            nextQueue.push(e.to)
          }
        }
      }
      queue = nextQueue
    }
  }

  private _init(s: number, t: number): number {
    const n = this._n
    this._exFlow = this._capSum < 2 ** 32 ? new Uint32Array(n) : Array(n).fill(0)
    this._potential[s] = n + 1
    this._bfs(t)
    const nexts = this._graph[s]
    for (let i = 0; i < nexts.length; i++) {
      const e = nexts[i]
      if (this._potential[e.to] < n) {
        this._graph[e.to][e.rev].cap = e.cap
        this._exFlow[s] -= e.cap
        this._exFlow[e.to] += e.cap
      }
      e.cap = 0
    }
    return this._calcActive(t)
  }

  private _push(u: number, t: number, e: Edge): boolean {
    const f = Math.min(e.cap, this._exFlow[u])
    const v = e.to
    e.cap -= f
    this._exFlow[u] -= f
    this._graph[v][e.rev].cap += f
    this._exFlow[v] += f
    if (this._exFlow[v] === f && v !== t) this._activeVertex.push(this._potential[v], v)
    return this._exFlow[u] === 0
  }

  private _discharge(u: number, t: number): number {
    const nexts = this._graph[u]
    for (let i = this._curEdge[u]; i < nexts.length; i++) {
      const e = nexts[i]
      if (this._potential[u] === this._potential[e.to] + 1 && e.cap > 0) {
        if (this._push(u, t, e)) return this._potential[u]
      }
    }
    return this._relabel(u)
  }

  private _globalRelabel(t: number): number {
    this._bfs(t)
    this._allVertex.clear()
    this._activeVertex.clear()
    return this._calcActive(t)
  }

  private _gapRelabel(u: number): void {
    const n = this._n
    for (let i = this._potential[u]; i <= this._height; i++) {
      for (let id = this._allVertex.data[n + i].next; id < n; id = this._allVertex.data[id].next) {
        this._potential[id] = n + 1
      }
      // eslint-disable-next-line no-multi-assign
      this._allVertex.data[n + i].prev = this._allVertex.data[n + i].next = n + i
    }
  }

  private _relabel(u: number): number {
    this._relabels++
    const prv = this._potential[u]
    let cur = this._n
    const nexts = this._graph[u]
    for (let i = 0; i < nexts.length; i++) {
      const e = nexts[i]
      if (cur > this._potential[e.to] + 1 && e.cap > 0) {
        this._curEdge[u] = i
        cur = this._potential[e.to] + 1
      }
    }
    if (this._allVertex.moreOne(prv)) {
      this._allVertex.erase(u)
      // eslint-disable-next-line no-cond-assign
      if ((this._potential[u] = cur) === this._n) {
        // eslint-disable-next-line no-return-assign
        return (this._potential[u] = this._n + 1), prv
      }
      this._activeVertex.push(cur, u)
      this._allVertex.insert(cur, u)
      this._height = Math.max(this._height, cur)
    } else {
      this._gapRelabel(u)
      // eslint-disable-next-line no-return-assign
      return (this._height = prv - 1)
    }
    return cur
  }
}

class _Stack {
  private readonly _n: number
  private readonly _h: number
  private readonly _node: Uint32Array

  constructor(n: number, h: number) {
    this._n = n
    this._h = h
    this._node = new Uint32Array(n + h)
    this.clear()
  }

  empty(h: number): boolean {
    return this._node[this._n + h] === this._n + h
  }

  top(h: number): number {
    return this._node[this._n + h]
  }

  pop(h: number): void {
    this._node[this._n + h] = this._node[this._node[this._n + h]]
  }

  push(h: number, u: number): void {
    this._node[u] = this._node[this._n + h]
    this._node[this._n + h] = u
  }

  clear(): void {
    for (let i = this._n; i < this._n + this._h; i++) this._node[i] = i
  }
}

class _List {
  private readonly _n: number
  private readonly _h: number
  readonly data: { prev: number; next: number }[]

  constructor(n: number, h: number) {
    this._n = n
    this._h = h
    this.data = Array(n + h)
    for (let i = 0; i < n + h; i++) {
      this.data[i] = {} as any
    }
    this.clear()
  }

  empty(h: number): boolean {
    return this.data[this._n + h].next === this._n + h
  }

  moreOne(h: number): boolean {
    return this.data[this._n + h].prev !== this.data[this._n + h].next
  }

  insert(h: number, u: number): void {
    this.data[u].prev = this.data[this._n + h].prev
    this.data[u].next = this._n + h
    this.data[this.data[this._n + h].prev].next = u
    this.data[this._n + h].prev = u
  }

  erase(u: number): void {
    this.data[this.data[u].prev].next = this.data[u].next
    this.data[this.data[u].next].prev = this.data[u].prev
  }

  clear(): void {
    for (let i = this._n; i < this._n + this._h; i++) {
      // eslint-disable-next-line no-multi-assign
      this.data[i].prev = this.data[i].next = i
    }
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
