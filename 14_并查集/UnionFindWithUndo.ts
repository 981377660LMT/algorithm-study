/** 可撤销并查集. */
class UnionFindArrayWithUndo {
  private readonly _n: number
  private readonly _parent: Uint32Array
  private readonly _rank: Uint32Array
  private readonly _optStack: { small: number; big: number; smallRank: number }[] = []
  private _part: number

  constructor(n: number) {
    this._n = n
    this._parent = new Uint32Array(n)
    this._rank = new Uint32Array(n)
    this._part = n
    for (let i = 0; i < n; i++) {
      this._parent[i] = i
      this._rank[i] = 1
    }
  }

  find(x: number): number {
    while (this._parent[x] !== x) {
      x = this._parent[x]
    }
    return x
  }

  union(x: number, y: number): boolean {
    let rootX = this.find(x)
    let rootY = this.find(y)
    if (rootX === rootY) {
      this._optStack.push({ small: -1, big: -1, smallRank: -1 })
      return false
    }
    if (this._rank[rootX] > this._rank[rootY]) {
      const tmp = rootX
      rootX = rootY
      rootY = tmp
    }
    this._parent[rootX] = rootY
    this._rank[rootY] += this._rank[rootX]
    this._part--
    this._optStack.push({ small: rootX, big: rootY, smallRank: this._rank[rootX] })
    return true
  }

  undo(): void {
    if (!this._optStack.length) {
      return
    }
    const { small: rootX, big: rootY, smallRank: rankX } = this._optStack.pop()!
    if (rootX === -1) {
      return
    }
    this._parent[rootX] = rootX
    this._rank[rootY] -= rankX
    this._part++
  }

  reset(): void {
    while (this._optStack.length) {
      this.undo()
    }
  }

  getState(): number {
    return this._optStack.length
  }

  rollback(state: number): boolean {
    if (state < 0 || state > this._optStack.length) return false
    while (this._optStack.length > state) {
      this.undo()
    }
    return true
  }

  isConnected(x: number, y: number): boolean {
    return this.find(x) === this.find(y)
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

  getSize(x: number): number {
    return this._rank[this.find(x)]
  }

  get part(): number {
    return this._part
  }
}

const _pool = new Map<unknown, number>()
function id(o: unknown): number {
  if (!_pool.has(o)) {
    _pool.set(o, _pool.size)
  }
  return _pool.get(o)!
}

class UnionFindMapWithUndo {
  private readonly _parent: Map<number, number> = new Map()
  private readonly _rank: Map<number, number> = new Map()
  private readonly _optStack: { small: number; big: number; smallRank: number }[] = []
  private _part = 0

  constructor(iterable?: Iterable<number>) {
    if (iterable) {
      for (const key of iterable) {
        this.add(key)
      }
    }
  }

  find(key: number): number {
    if (!this._parent.has(key)) {
      this.add(key)
      return key
    }
    while ((this._parent.get(key) || key) !== key) {
      key = this._parent.get(key)!
    }
    return key
  }

  union(key1: number, key2: number): boolean {
    let root1 = this.find(key1)
    let root2 = this.find(key2)
    if (root1 === root2) {
      this._optStack.push({ small: -1, big: -1, smallRank: -1 })
      return false
    }
    if (this._rank.get(root1)! > this._rank.get(root2)!) {
      const tmp = root1
      root1 = root2
      root2 = tmp
    }
    this._parent.set(root1, root2)
    this._rank.set(root2, this._rank.get(root2)! + this._rank.get(root1)!)
    this._part--
    this._optStack.push({ small: root1, big: root2, smallRank: this._rank.get(root1)! })
    return true
  }

  undo(): void {
    if (!this._optStack.length) {
      return
    }
    const { small: root1, big: root2, smallRank: rank1 } = this._optStack.pop()!
    if (root1 === -1) {
      return
    }
    this._parent.set(root1, root1)
    this._rank.set(root2, this._rank.get(root2)! - rank1)
    this._part++
  }

  reset(): void {
    while (this._optStack.length) {
      this.undo()
    }
  }

  isConnected(key1: number, key2: number): boolean {
    return this.find(key1) === this.find(key2)
  }

  getGroups(): Map<number, number[]> {
    const groups = new Map<number, number[]>()
    for (const key of this._parent.keys()) {
      const root = this.find(key)
      if (!groups.has(root)) {
        groups.set(root, [])
      }
      groups.get(root)!.push(key)
    }
    return groups
  }

  add(key: number): boolean {
    if (this._parent.has(key)) {
      return false
    }
    this._parent.set(key, key)
    this._rank.set(key, 1)
    this._part++
    return true
  }

  has(key: number): boolean {
    return this._parent.has(key)
  }

  getSize(x: number): number {
    return this._rank.get(this.find(x)) || 0
  }

  get part(): number {
    return this._part
  }
}

/** @deprecated */
class UnionFindArrayWithTimeTravel {
  private readonly _n: number
  private readonly _data: Int32Array
  private readonly _history: { root: number; data: number }[] = []
  private _part = 0
  private _innerSnap = 0

  constructor(n: number) {
    this._n = n
    this._data = new Int32Array(n).fill(-1)
    this._part = n
  }

  /** 撤销上一次合并操作(包含未合并成功的操作). */
  undo(): boolean {
    if (!this._history.length) return false
    const { root: small, data: smallData } = this._history.pop()!
    const { root: big, data: bigData } = this._history.pop()!
    this._data[small] = smallData
    this._data[big] = bigData
    this._part += +(big !== small)
    return true
  }

  /**
   * 保存并查集当前的状态.
   * `Snapshot()` 之后可以调用 `Rollback(-1)` 回滚到这个状态.
   */
  snapShot(): void {
    this._innerSnap = this._history.length >>> 1
  }

  /**
   * 回滚到指定的状态.
   * `-1` 表示回滚到上一次 `SnapShot` 时保存的状态.
   * 其他值表示回滚到状态为`toState`时的状态.
   */
  rollback(toState = -1): boolean {
    if (toState === -1) toState = this._innerSnap
    toState <<= 1
    if (toState < 0 || toState > this._history.length) {
      return false
    }
    while (toState < this._history.length) {
      this.undo()
    }
    return true
  }

  /**
   * 获取当前并查集的状态id.
   * 也即当前合并`union`被调用的次数.
   */
  getState(): number {
    return this._history.length >>> 1
  }

  reset(): void {
    while (this._history.length) {
      this.undo()
    }
  }

  /**
   * 按秩合并.
   */
  union(u: number, v: number, f?: (big: number, small: number) => void): boolean {
    u = this.find(u)
    v = this.find(v)
    this._history.push({ root: u, data: this._data[u] })
    this._history.push({ root: v, data: this._data[v] })
    if (u === v) return false
    if (this._data[u] > this._data[v]) {
      const tmp = u
      u = v
      v = tmp
    }
    this._data[u] += this._data[v]
    this._data[v] = u
    this._part--
    f && f(u, v)
    return true
  }

  /** 因为需要支持撤销，所以不进行路径压缩. */
  find(x: number): number {
    let cur = x
    while (this._data[cur] >= 0) {
      cur = this._data[cur]
    }
    return cur
  }

  isConnected(u: number, v: number): boolean {
    return this.find(u) === this.find(v)
  }

  getSize(u: number): number {
    return -this._data[this.find(u)]
  }

  getGroups(): Map<number, number[]> {
    const res = new Map<number, number[]>()
    for (let i = 0; i < this._n; i++) {
      const root = this.find(i)
      !res.has(root) && res.set(root, [])
      res.get(root)!.push(i)
    }
    return res
  }

  get part(): number {
    return this._part
  }
}

export { UnionFindArrayWithUndo, UnionFindMapWithUndo }

if (require.main === module) {
  const uf = new UnionFindArrayWithUndo(5)
  console.log(uf.getGroups())
  uf.union(1, 3)
  console.log(uf.getGroups())
  uf.undo()
  console.log(uf.getGroups())
}
