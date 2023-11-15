// 无向图中：
// 联通分量数(part) = 树的个数(treeCount) + 环的个数
// 边的个数: 总点数 - 树的个数
// 树的性质: 联通分量中的点数 = 边数 + 1
// 环的性质: 联通分量中的点数 = 边数

class UnionFindGraphArray {
  private _part: number
  private _treeCount: number
  private readonly _n: number
  private readonly _data: Int32Array
  private readonly _edge: Uint32Array
  private readonly _history: { root: number; data: number; edge: number }[] = []

  constructor(n: number) {
    this._part = n
    this._treeCount = n
    this._n = n
    this._data = new Int32Array(n).fill(-1)
    this._edge = new Uint32Array(n)
  }

  /**
   * 添加边对(u,v).
   * @alias addPair
   */
  union(u: number, v: number): boolean {
    u = this.find(u)
    v = this.find(v)
    this._history.push({ root: u, data: this._data[u], edge: this._edge[u] }) // big
    this._history.push({ root: v, data: this._data[v], edge: this._edge[v] }) // small
    if (u === v) {
      if (this.isTree(u)) {
        this._treeCount--
      }
      this._edge[u]++
      return false
    }

    if (this._data[u] > this._data[v]) {
      u ^= v
      v ^= u
      u ^= v
    }
    if (this.isTree(u) || this.isTree(v)) {
      this._treeCount--
    }
    this._data[u] += this._data[v]
    this._data[v] = u
    this._edge[u] += this._edge[v] + 1
    this._part--
    return true
  }

  /**
   * 不能路径压缩.
   */
  find(u: number): number {
    let cur = u
    while (this._data[cur] >= 0) {
      cur = this._data[cur]
    }
    return cur
  }

  /**
   * 撤销上一次合并操作，没合并成功也要撤销.
   */
  undo(): boolean {
    if (!this._history.length) {
      return false
    }
    const { root: small, data: smallData, edge: smallEdge } = this._history.pop()!
    const { root: big, data: bigData, edge: bigEdge } = this._history.pop()!
    this._data[small] = smallData
    this._data[big] = bigData
    this._edge[small] = smallEdge
    this._edge[big] = bigEdge
    if (big === small) {
      if (this.isTree(big)) {
        this._treeCount++
      }
    } else {
      if (this.isTree(big) || this.isTree(small)) {
        this._treeCount++
      }
      this._part++
    }
    return true
  }

  /**
   * 从每条边中恰好选一个点, 最多能选出多少个不同的点.
   * 对每个大小为m的连通块,树的贡献为m-1,环的贡献为m.
   * 因此答案为`总点数-树的个数`.
   */
  solve(): number {
    return this._n - this._treeCount
  }

  /**
   * u所在的联通分量是否为树.
   */
  isTree(u: number): boolean {
    const root = this.find(u)
    const vertex = this.getSize(root)
    return vertex === this._edge[root] + 1
  }

  /**
   * u所在的联通分量是否为环.
   */
  isCycle(u: number): boolean {
    const root = this.find(u)
    const vertex = this.getSize(root)
    return vertex === this._edge[root]
  }

  /**
   * 联通分量为树的联通分量个数(孤立点也算树).
   */
  countTree(): number {
    return this._treeCount
  }

  /**
   * 联通分量为环的联通分量个数(孤立点不算环).
   */
  countCycle(): number {
    return this._part - this._treeCount
  }

  countEdge(): number {
    return this._n - this._treeCount
  }

  getSize(x: number): number {
    return -this._data[this.find(x)]
  }

  getEdge(x: number): number {
    return this._edge[this.find(x)]
  }

  getGroups(): Map<number, number[]> {
    const groups = new Map<number, number[]>()
    for (let i = 0; i < this._n; i++) {
      const root = this.find(i)
      if (!groups.has(root)) {
        groups.set(root, [])
      }
      groups.get(root)!.push(i)
    }
    return groups
  }

  get part(): number {
    return this._part
  }

  get treeCount(): number {
    return this._treeCount
  }
}

class UnionFindGraphMap {
  private _part = 0
  private _treeCount = 0
  private readonly _data = new Map<number, number>()
  private readonly _edge = new Map<number, number>()
  private readonly _history: { root: number; data: number; edge: number }[] = []

  /**
   * 添加边对(u,v).
   * @alias addPair
   */
  union(u: number, v: number): boolean {
    u = this.find(u)
    v = this.find(v)
    this._history.push({ root: u, data: this._data.get(u)!, edge: this._edge.get(u)! }) // big
    this._history.push({ root: v, data: this._data.get(v)!, edge: this._edge.get(v)! }) // small
    if (u === v) {
      if (this.isTree(u)) {
        this._treeCount--
      }
      this._edge.set(u, this._edge.get(u)! + 1)
      return false
    }

    if (this._data.get(u)! > this._data.get(v)!) {
      u ^= v
      v ^= u
      u ^= v
    }
    if (this.isTree(u) || this.isTree(v)) {
      this._treeCount--
    }
    this._data.set(u, this._data.get(u)! + this._data.get(v)!)
    this._data.set(v, u)
    this._edge.set(u, this._edge.get(u)! + this._edge.get(v)! + 1)
    this._part--
    return true
  }

  /**
   * 不能路径压缩.
   */
  find(u: number): number {
    if (!this._data.has(u)) {
      this.add(u)
      return u
    }
    let cur = u
    while ((this._data.get(cur) || -1) >= 0) {
      cur = this._data.get(cur)!
    }
    return cur
  }

  /**
   * 撤销上一次合并操作，没合并成功也要撤销.
   */
  undo(): boolean {
    if (!this._history.length) {
      return false
    }
    const { root: small, data: smallData, edge: smallEdge } = this._history.pop()!
    const { root: big, data: bigData, edge: bigEdge } = this._history.pop()!
    this._data.set(small, smallData)
    this._data.set(big, bigData)
    this._edge.set(small, smallEdge)
    this._edge.set(big, bigEdge)
    if (big === small) {
      if (this.isTree(big)) {
        this._treeCount++
      }
    } else {
      if (this.isTree(big) || this.isTree(small)) {
        this._treeCount++
      }
      this._part++
    }
    return true
  }

  /**
   * 从每条边中恰好选一个点, 最多能选出多少个不同的点.
   * 对每个大小为m的连通块,树的贡献为m-1,环的贡献为m.
   * 因此答案为`总点数-树的个数`.
   */
  solve(): number {
    return this._data.size - this._treeCount
  }

  /**
   * u所在的联通分量是否为树.
   */
  isTree(u: number): boolean {
    const root = this.find(u)
    const vertex = this.getSize(root)
    return vertex === this._edge.get(root)! + 1
  }

  /**
   * u所在的联通分量是否为环.
   */
  isCycle(u: number): boolean {
    const root = this.find(u)
    const vertex = this.getSize(root)
    return vertex === this._edge.get(root)!
  }

  /**
   * 联通分量为树的联通分量个数(孤立点也算树).
   */
  countTree(): number {
    return this._treeCount
  }

  /**
   * 联通分量为环的联通分量个数(孤立点不算环).
   */
  countCycle(): number {
    return this._part - this._treeCount
  }

  countEdge(): number {
    return this._data.size - this._treeCount
  }

  getSize(x: number): number {
    return -this._data.get(this.find(x))!
  }

  getEdge(x: number): number {
    return this._edge.get(this.find(x))!
  }

  getGroups(): Map<number, number[]> {
    const groups = new Map<number, number[]>()
    for (const k of this._data.keys()) {
      const root = this.find(k)
      if (!groups.has(root)) {
        groups.set(root, [])
      }
      groups.get(root)!.push(k)
    }
    return groups
  }

  add(u: number): boolean {
    if (this._data.has(u)) {
      return false
    }
    this._data.set(u, -1)
    this._edge.set(u, 0)
    this._part++
    this._treeCount++
    return true
  }

  get part(): number {
    return this._part
  }

  get treeCount(): number {
    return this._treeCount
  }
}

export { UnionFindGraphArray, UnionFindGraphMap }
