function minOperationsQueries(n: number, edges: number[][], queries: number[][]): number[] {
  const tree = new Tree(n)
  const weights: Map<number, number>[] = Array(n)
  for (let i = 0; i < n; i++) weights[i] = new Map()
  edges.forEach(([u, v, w]) => {
    tree.addEdge(u, v, w)
    weights[u].set(v, w)
    weights[v].set(u, w)
  })
  tree.build(0)

  // 对每个查询求出边权，答案为边数减去最大频率
  const res = Array(queries.length).fill(0)
  for (let i = 0; i < queries.length; i++) {
    const [u, v] = queries[i]
    const path = tree.getPath(u, v)
    const weightCounter = new Uint16Array(27)
    for (let j = 0; j < path.length - 1; j++) {
      const pre = path[j]
      const cur = path[j + 1]
      const w = weights[pre].get(cur)!
      weightCounter[w]++
    }
    res[i] = path.length - 1 - Math.max(...weightCounter)
  }

  return res
}

class Tree {
  readonly depth: Uint32Array
  readonly parent: Int32Array
  readonly depthWeighted: number[]
  readonly lid: Uint32Array
  readonly rid: Uint32Array
  private readonly _idToNode: Uint32Array
  private readonly _top: Uint32Array
  private readonly _heavySon: Int32Array
  private _timer = 0

  /** 链式前向星存图 */
  private readonly _n: number
  private readonly _preEdge: Int32Array
  private readonly _lastEdge: Int32Array
  private readonly _edgeTo: Int32Array
  private readonly _weight: number[]
  private _edgeId = 0

  constructor(n: number) {
    this.parent = new Int32Array(n)
    this.depth = new Uint32Array(n)
    this.lid = new Uint32Array(n)
    this.rid = new Uint32Array(n)
    this.depthWeighted = Array(n)
    this._idToNode = new Uint32Array(n)
    this._top = new Uint32Array(n)
    this._heavySon = new Int32Array(n)
    this._lastEdge = new Int32Array(n)
    this._preEdge = new Int32Array(2 * (n - 1))
    this._edgeTo = new Int32Array(2 * (n - 1))
    this._weight = Array(2 * (n - 1))
    this._n = n
    for (let i = 0; i < n; i++) {
      this.parent[i] = -1
      this.depthWeighted[i] = 0
      this._lastEdge[i] = -1
    }
    for (let i = 0; i < 2 * (n - 1); i++) {
      this._preEdge[i] = -1
      this._edgeTo[i] = -1
      this._weight[i] = 0
    }
  }

  addEdge(from: number, to: number, weight = 1): void {
    this.addDirectedEdge(from, to, weight)
    this.addDirectedEdge(to, from, weight)
  }

  addDirectedEdge(from: number, to: number, weight = 1): void {
    const eid = this._edgeId++
    this._preEdge[eid] = this._lastEdge[from]
    this._lastEdge[from] = eid
    this._edgeTo[eid] = to
    this._weight[eid] = weight
  }

  /**
   * 当 root 为 `-1(默认值)`时, 会从`0`开始遍历未访问过的连通分量.
   */
  build(root = -1): void {
    if (root === -1) {
      for (let i = 0; i < this._n; i++) {
        if (this.parent[i] === -1) {
          this._build(i, -1, 0, 0)
          this._markTop(i, i)
        }
      }
      return
    }

    this._build(root, -1, 0, 0)
    this._markTop(root, root)
  }

  /**
   * 返回 root 的欧拉序区间, 左闭右开, 0-indexed.
   */
  id(root: number): [inId: number, outId: number] {
    return [this.lid[root], this.rid[root]]
  }

  /**
   * 返回边 u-v 对应的 欧拉序起点编号, 0-indexed.
   */
  eid(u: number, v: number): number {
    const id1 = this.lid[u]
    const id2 = this.lid[v]
    return id1 > id2 ? id1 : id2
  }

  lca(u: number, v: number): number {
    while (true) {
      if (this.lid[u] > this.lid[v]) {
        u ^= v
        v ^= u
        u ^= v
      }
      if (this._top[u] === this._top[v]) {
        return u
      }
      v = this.parent[this._top[v]]
    }
  }

  /**
   * 返回 root 为根时, u 和 v 的最近公共祖先.
   */
  rootedLca(u: number, v: number, root: number): number {
    return this.lca(u, v) ^ this.lca(u, root) ^ this.lca(v, root)
  }

  rootedParent(u: number, root: number): number {
    return this.jump(u, root, 1)
  }

  dist(u: number, v: number, weighted: boolean): number {
    if (weighted) {
      return this.depthWeighted[u] + this.depthWeighted[v] - 2 * this.depthWeighted[this.lca(u, v)]
    }
    return this.depth[u] + this.depth[v] - 2 * this.depth[this.lca(u, v)]
  }

  /**
   * 返回 root 的第 k 个祖先, k 从 0 开始计数.
   * 如果不存在这样的祖先, 返回-1.
   */
  kthAncestor(root: number, k: number): number {
    if (k > this.depth[root]) {
      return -1
    }
    while (true) {
      const u = this._top[root]
      if (this.lid[root] - k >= this.lid[u]) {
        return this._idToNode[this.lid[root] - k]
      }
      k -= this.lid[root] - this.lid[u] + 1
      root = this.parent[u]
    }
  }

  /**
   * 从 from 节点跳向 to 节点,跳过 step 个节点(0-indexed), 返回跳到的节点.
   * 如果不存在这样的节点, 返回-1.
   */
  jump(from: number, to: number, step: number): number {
    if (step === 1) {
      if (from === to) {
        return -1
      }
      if (this.isInSubtree(to, from)) {
        return this.kthAncestor(to, this.depth[to] - this.depth[from] - 1)
      }
      return this.parent[from]
    }

    const lca = this.lca(from, to)
    const dac = this.depth[from] - this.depth[lca]
    const dbc = this.depth[to] - this.depth[lca]
    if (step > dac + dbc) {
      return -1
    }
    if (step <= dac) {
      return this.kthAncestor(from, step)
    }
    return this.kthAncestor(to, dac + dbc - step)
  }

  collectChildren(root: number): number[] {
    const res: number[] = []
    for (let eid = this._lastEdge[root]; ~eid; eid = this._preEdge[eid]) {
      const next = this._edgeTo[eid]
      if (next !== this.parent[root]) {
        res.push(next)
      }
    }
    return res
  }

  /**
   * 返回沿着`路径顺序`的 [起点,终点] 的 欧拉序 `左闭右闭` 数组.
   * !eg:[[2 0] [4 4]] 沿着路径顺序但不一定沿着欧拉序.
   */
  getPathDecomposition(u: number, v: number, vertex: boolean): [from: number, to: number][] {
    const up: [start: number, end: number][] = []
    const down: [start: number, end: number][] = []
    while (true) {
      if (this._top[u] === this._top[v]) {
        break
      }

      if (this.lid[u] < this.lid[v]) {
        down.push([this.lid[this._top[v]], this.lid[v]])
        v = this.parent[this._top[v]]
      } else {
        up.push([this.lid[u], this.lid[this._top[u]]])
        u = this.parent[this._top[u]]
      }
    }
    const offset = vertex ? 0 : 1
    if (this.lid[u] < this.lid[v]) {
      down.push([this.lid[u] + offset, this.lid[v]])
    } else if (this.lid[v] + offset <= this.lid[u]) {
      up.push([this.lid[u], this.lid[v] + offset])
    }
    up.push(...down.reverse())
    return up
  }

  /**
   * 遍历路径上的 `[起点,终点)` 欧拉序 `左闭右开` 区间.
   */
  enumeratePathDecomposition(
    u: number,
    v: number,
    vertex: boolean,
    callback: (start: number, end: number) => void
  ): void {
    while (true) {
      if (this._top[u] === this._top[v]) {
        break
      }

      if (this.lid[u] < this.lid[v]) {
        const a = this.lid[this._top[v]]
        const b = this.lid[v]
        a < b ? callback(a, b + 1) : callback(b, a + 1)
        v = this.parent[this._top[v]]
      } else {
        const a = this.lid[u]
        const b = this.lid[this._top[u]]
        a < b ? callback(a, b + 1) : callback(b, a + 1)
        u = this.parent[this._top[u]]
      }
    }

    const offset = vertex ? 0 : 1
    if (this.lid[u] < this.lid[v]) {
      const a = this.lid[u] + offset
      const b = this.lid[v]
      a < b ? callback(a, b + 1) : callback(b, a + 1)
    } else if (this.lid[v] + offset <= this.lid[u]) {
      const a = this.lid[u]
      const b = this.lid[v] + offset
      a < b ? callback(a, b + 1) : callback(b, a + 1)
    }
  }

  getPath(u: number, v: number): number[] {
    const res: number[] = []
    const composition = this.getPathDecomposition(u, v, true)
    composition.forEach(([start, end]) => {
      if (start <= end) {
        for (let i = start; i <= end; i++) {
          res.push(this._idToNode[i])
        }
      } else {
        for (let i = start; i >= end; i--) {
          res.push(this._idToNode[i])
        }
      }
    })
    return res
  }

  /**
   * 以root为根时,结点v的子树大小.
   */
  subSize(v: number, root = -1): number {
    if (root === -1) {
      return this.rid[v] - this.lid[v]
    }
    if (v === root) {
      return this._n
    }
    const x = this.jump(v, root, 1)
    if (this.isInSubtree(v, x)) {
      return this.rid[v] - this.lid[v]
    }
    return this._n - this.rid[x] + this.lid[x]
  }

  /**
   * child 是否在 root 的子树中 (child和root不能相等).
   */
  isInSubtree(child: number, root: number): boolean {
    return this.lid[root] <= this.lid[child] && this.rid[child] <= this.rid[root]
  }

  /**
   * 寻找以 start 为 top 的重链 , heavyPath[-1] 即为重链底端节点.
   */
  getHeavyPath(start: number): number[] {
    const res: number[] = [start]
    let cur = start
    while (~this._heavySon[cur]) {
      cur = this._heavySon[cur]
      res.push(cur)
    }
    return res
  }

  /**
   * 返回结点v的重儿子.如果没有重儿子,返回-1.
   */
  getHeavyChild(v: number): number {
    const k = this.lid[v] + 1
    if (k === this._n) {
      return -1
    }
    const w = this._idToNode[k]
    if (this.parent[w] === v) {
      return w
    }
    return -1
  }

  toAdjList(): [next: number, weight: number, edgeId: number][][] {
    const res = Array(this._n)
    for (let cur = 0; cur < this._n; cur++) {
      const nexts: [next: number, weight: number, edgeId: number][] = []
      for (let eid = this._lastEdge[cur]; ~eid; eid = this._preEdge[eid]) {
        const next = this._edgeTo[eid]
        const weight = this._weight[eid]
        nexts.push([next, weight, eid])
      }
      res[cur] = nexts.reverse()
    }
    return res
  }

  private _build(cur: number, pre: number, dep: number, dist: number): number {
    let subSize = 1
    let heavySon = -1
    let heavySize = 0
    for (let eid = this._lastEdge[cur]; ~eid; eid = this._preEdge[eid]) {
      const next = this._edgeTo[eid]
      if (next !== pre) {
        const nextSize = this._build(next, cur, dep + 1, dist + this._weight[eid])
        subSize += nextSize
        if (nextSize > heavySize) {
          heavySize = nextSize
          heavySon = next
        }
      }
    }
    this.depth[cur] = dep
    this.parent[cur] = pre
    this.depthWeighted[cur] = dist
    this._heavySon[cur] = heavySon
    return subSize
  }

  private _markTop(cur: number, top: number): void {
    this._top[cur] = top
    this.lid[cur] = this._timer
    this._idToNode[this._timer] = cur
    this._timer++
    const heavySon = this._heavySon[cur]
    if (~heavySon) {
      this._markTop(heavySon, top)
      for (let eid = this._lastEdge[cur]; ~eid; eid = this._preEdge[eid]) {
        const next = this._edgeTo[eid]
        if (next !== heavySon && next !== this.parent[cur]) {
          this._markTop(next, next)
        }
      }
    }
    this.rid[cur] = this._timer
  }
}

export {}