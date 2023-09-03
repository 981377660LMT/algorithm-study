/* eslint-disable no-inner-declarations */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

/**
 * 二分图最大匹配.
 */
class BipartiteMathing {
  private readonly _graph: number[][]
  private readonly _alive: Uint8Array
  private readonly _used: Uint32Array
  private readonly _match: Int32Array
  private readonly _n: number
  private _timestamp = 0

  constructor(n: number) {
    const graph: number[][] = Array(n)
    for (let i = 0; i < n; i++) graph[i] = []
    const alive = new Uint8Array(n)
    const used = new Uint32Array(n)
    const match = new Int32Array(n)
    for (let i = 0; i < n; i++) {
      alive[i] = 1
      match[i] = -1
    }
    this._graph = graph
    this._alive = alive
    this._used = used
    this._match = match
    this._n = n
  }

  /**
   * 添加边 `left - right`.
   * @param left 左侧顶点.0 <= left < L.
   * @param right 右侧顶点.L <= right < n.
   */
  addEdge(left: number, right: number): void {
    this._graph[left].push(right)
    this._graph[right].push(left)
  }

  removeEdge(u: number, v: number): void {
    const nextsU = this._graph[u]
    for (let i = 0; i < nextsU.length; i++) {
      if (nextsU[i] === v) {
        nextsU.splice(i, 1)
        break
      }
    }
    const nextsV = this._graph[v]
    for (let i = 0; i < nextsV.length; i++) {
      if (nextsV[i] === u) {
        nextsV.splice(i, 1)
        break
      }
    }
  }

  /**
   * `O(VE)`求二分图最大匹配.
   */
  maxMatching(): [left: number, right: number][] {
    this._match.fill(-1) // 重置匹配
    for (let u = 0; u < this._n; u++) {
      if (this._alive[u] && this._match[u] === -1) {
        this._timestamp++
        this._argument(u)
      }
    }

    const res: [left: number, right: number][] = []
    for (let u = 0; u < this._n; u++) {
      const cand = this._match[u]
      if (u < cand) res.push([u, cand])
    }
    return res
  }

  /**
   * 删除顶点 `idx`, 返回流量的变化量(-1/0).
   */
  removeVertex(idx: number): -1 | 0 {
    this._alive[idx] = 0
    if (this._match[idx] === -1) return 0
    this._match[this._match[idx]] = -1
    this._timestamp++
    const res = this._argument(this._match[idx])
    this._match[idx] = -1
    return res ? 0 : -1
  }

  /**
   * 添加顶点 `idx`, 返回流量的变化量(0/1).
   */
  addVertex(idx: number): 0 | 1 {
    this._alive[idx] = 1
    this._timestamp++
    return this._argument(idx) ? 1 : 0
  }

  /**
   * 获取匹配边.需要先调用 `maxMatching`.
   */
  getMatchingEdges(): [left: number, right: number][] {
    const res: [left: number, right: number][] = []
    for (let u = 0; u < this._n; u++) {
      const cand = this._match[u]
      if (u < cand) res.push([u, cand])
    }
    return res
  }

  private _argument(idx: number): boolean {
    this._used[idx] = this._timestamp
    const nexts = this._graph[idx]
    for (let i = 0; i < nexts.length; i++) {
      const to = nexts[i]
      const toMatch = this._match[to]
      if (!this._alive[to]) continue
      if (toMatch === -1 || (this._used[toMatch] !== this._timestamp && this._argument(toMatch))) {
        this._match[idx] = to
        this._match[to] = idx
        return true
      }
    }
    return false
  }
}

/**
 * 从边创建二分图最大匹配.
 * @param n 顶点数.
 * @param edges 边集.
 * @returns
 * bm: 二分图最大匹配
 * ids: 原图中的点在二分图中的编号.左侧点编号为0-L-1, 右侧点编号为L-n-1.
 * rids: 二分图中的点在原图中的编号.0-n-1 -> 0-n-1.
 * @example
 * ```ts
 * const [bm, ids, rids] = createBipartiteMathingFromEdges(4, [[0, 1], [1, 2], [2, 3]])
 * const M = bm.maxMatching()
 * const edges = M.map(([u, v]) => [rids[u], rids[v]])
 * console.log(edges) // [[0, 1], [2, 3]]
 * ```
 */
function createBipartiteMathingFromEdges(
  n: number,
  edges: [u: number, v: number][]
): [bm: BipartiteMathing, ids: Uint32Array, rids: Uint32Array] {
  const graph: number[][] = Array(n)
  for (let i = 0; i < n; i++) graph[i] = []
  edges.forEach(([u, v]) => {
    graph[u].push(v)
    graph[v].push(u)
  })

  const [colors, ok] = isBipartite(n, graph)
  if (!ok) throw new Error('not bipartite')

  let leftCount = 0
  for (let i = 0; i < n; i++) leftCount += +!colors[i] // 规定左侧点颜色为0, 右侧点颜色为1
  const ids = new Uint32Array(n)
  const rids = new Uint32Array(n)
  let left = 0
  let right = 0
  for (let i = 0; i < n; i++) {
    if (!colors[i]) {
      ids[i] = left
      rids[left] = i
      left++
    } else {
      ids[i] = right + leftCount
      rids[right + leftCount] = i
      right++
    }
  }

  const bm = new BipartiteMathing(n)
  edges.forEach(([u, v]) => {
    if (colors[u]) {
      u ^= v
      v ^= u
      u ^= v
    }
    bm.addEdge(ids[u], ids[v])
  })

  return [bm, ids, rids]
}

/**
 * 判断是否是二分图.返回 (染色的01数组, 是否是二分图).
 */
function isBipartite(n: number, graph: number[][]): [colors: Int8Array, ok: boolean] {
  const colors = new Int8Array(n).fill(-1)
  for (let i = 0; i < n; i++) {
    if (colors[i] === -1) {
      colors[i] = 0
      const stack = [i]
      while (stack.length) {
        const cur = stack.pop()!
        const nexts = graph[cur]
        for (let j = 0; j < nexts.length; j++) {
          const next = nexts[j]
          if (colors[next] === -1) {
            colors[next] = 1 ^ colors[cur]
            stack.push(next)
          } else if (colors[next] === colors[cur]) {
            return [colors, false]
          }
        }
      }
    }
  }
  return [colors, true]
}

export { BipartiteMathing, createBipartiteMathingFromEdges, isBipartite }

if (require.main === module) {
  // // createBipartiteMathingFromEdges
  // const [bm, _, rids] = createBipartiteMathingFromEdges(4, [
  //   [0, 1],
  //   [1, 2],
  //   [2, 3]
  // ])
  // const M = bm.maxMatching()
  // const edges = M.map(([u, v]) => [rids[u], rids[v]])
  // console.log(edges) // [[0, 1], [2, 3]]

  // https://atcoder.jp/contests/abc317/tasks/abc317_g
  // 跑m次匈牙利，每跑一次就删去完美匹配的边
  function rearrange(grid: number[][]): number[][] {
    const ROW = grid.length
    const COL = grid[0].length
    const res: number[][] = Array(ROW)
    for (let i = 0; i < ROW; i++) res[i] = Array(COL).fill(0)

    const G = new BipartiteMathing(ROW + ROW)
    for (let i = 0; i < ROW; i++) {
      const row = grid[i]
      row.forEach(v => {
        G.addEdge(i, ROW + v)
      })
    }

    for (let c = 0; c < COL; c++) {
      const matching = G.maxMatching()
      matching.forEach(([u, v]) => {
        res[u][c] = v - ROW + 1
        G.removeEdge(u, v)
      })
    }

    return res
  }
}
