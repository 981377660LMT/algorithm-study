/* eslint-disable @typescript-eslint/no-non-null-assertion */

class UnionFindArrayWithUndo {
  private readonly _n: number
  private readonly _parent: Uint32Array
  private readonly _rank: Uint32Array
  private readonly _optStack: [small: number, big: number, smallRank: number][] = []
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
      this._optStack.push([-1, -1, -1])
      return false
    }
    if (this._rank[rootX] > this._rank[rootY]) {
      rootX ^= rootY
      rootY ^= rootX
      rootX ^= rootY
    }
    this._parent[rootX] = rootY
    this._rank[rootY] += this._rank[rootX]
    this._part--
    this._optStack.push([rootX, rootY, this._rank[rootX]])
    return true
  }

  undo(): void {
    if (!this._optStack.length) {
      return
    }
    const [rootX, rootY, rankX] = this._optStack.pop()!
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
  private readonly _optStack: [small: number, big: number, smallRank: number][] = []
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
      this._optStack.push([-1, -1, -1])
      return false
    }
    if (this._rank.get(root1)! > this._rank.get(root2)!) {
      root1 ^= root2
      root2 ^= root1
      root1 ^= root2
    }
    this._parent.set(root1, root2)
    this._rank.set(root2, this._rank.get(root2)! + this._rank.get(root1)!)
    this._part--
    this._optStack.push([root1, root2, this._rank.get(root1)!])
    return true
  }

  undo(): void {
    if (!this._optStack.length) {
      return
    }
    const [root1, root2, rank1] = this._optStack.pop()!
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

  get part(): number {
    return this._part
  }
}

export { UnionFindArrayWithUndo, UnionFindMapWithUndo }
