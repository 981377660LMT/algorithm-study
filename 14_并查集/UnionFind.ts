/* eslint-disable @typescript-eslint/no-non-null-assertion */

class UnionFindArray {
  private readonly _n: number
  private readonly _data: Int32Array
  private _part: number

  constructor(n: number) {
    this._n = n
    this._part = n
    this._data = new Int32Array(n).fill(-1)
  }

  union(x: number, y: number, beforeUnion?: (big: number, small: number) => void): boolean {
    let rootX = this.find(x)
    let rootY = this.find(y)
    if (rootX === rootY) return false
    if (this._data[rootX] > this._data[rootY]) {
      const tmp = rootX
      rootX = rootY
      rootY = tmp
    }
    beforeUnion && beforeUnion(rootX, rootY)
    this._data[rootX] += this._data[rootY]
    this._data[rootY] = rootX
    this._part -= 1
    return true
  }

  find(x: number): number {
    // eslint-disable-next-line no-return-assign
    return this._data[x] < 0 ? x : (this._data[x] = this.find(this._data[x]))
  }

  isConnected(x: number, y: number): boolean {
    return this.find(x) === this.find(y)
  }

  getGroups(): Map<number, number[]> {
    const groups = new Map<number, number[]>()
    for (let i = 0; i < this._n; i++) {
      const root = this.find(i)
      !groups.has(root) && groups.set(root, [])
      groups.get(root)!.push(i)
    }
    return groups
  }

  getRoots(): number[] {
    const res = Array(this._n)
    for (let i = 0; i < this._n; i++) {
      res[i] = this.find(i)
    }
    return [...new Set(res)]
  }

  getPart(): number {
    return this._part
  }

  getSize(x: number): number {
    return -this._data[this.find(x)]
  }

  clear(): void {
    this._part = this._n
    this._data.fill(-1)
  }

  toString(): string {
    return [...this.getGroups()].map(([root, member]) => `${root}: ${member}`).join('\n')
  }
}

class UnionFindMap<V extends number | string> {
  private readonly _parent: Map<V, V> = new Map()
  private readonly _rank: Map<V, number> = new Map()
  private readonly _autoAdd: boolean
  private _part = 0

  constructor(arrayLike: ArrayLike<V> = [], autoAdd = true) {
    this._autoAdd = autoAdd
    for (let i = 0; i < arrayLike.length; i++) {
      this.add(arrayLike[i])
    }
  }

  union(x: V, y: V, beforeUnion?: (big: V, small: V) => void): boolean {
    let rootX = this.find(x)
    let rootY = this.find(y)
    if (rootX === rootY) return false
    if (this._rank.get(rootX)! > this._rank.get(rootY)!) {
      ;[rootX, rootY] = [rootY, rootX]
    }
    beforeUnion && beforeUnion(rootY, rootX)
    this._parent.set(rootX, rootY)
    this._rank.set(rootY, this._rank.get(rootY)! + this._rank.get(rootX)!)
    this._part -= 1
    return true
  }

  find(x: V): V {
    if (!this._parent.has(x)) {
      if (this._autoAdd) {
        this.add(x)
      }
      return x
    }

    while ((this._parent.get(x) || x) !== x) {
      this._parent.set(x, this._parent.get(this._parent.get(x)!)!)
      x = this._parent.get(x)!
    }
    return x
  }

  add(x: V): boolean {
    if (this._parent.has(x)) {
      return false
    }
    this._parent.set(x, x)
    this._rank.set(x, 1)
    this._part += 1
    return true
  }

  isConnected(x: V, y: V): boolean {
    return this.find(x) === this.find(y)
  }

  getGroups(): Map<V, V[]> {
    const groups = new Map<V, V[]>()
    for (const key of this._parent.keys()) {
      const root = this.find(key)
      !groups.has(root) && groups.set(root, [])
      groups.get(root)!.push(key)
    }
    return groups
  }

  getRoots(): V[] {
    const res = new Set<V>()
    this._rank.forEach((_, key) => res.add(this.find(key)))
    return [...new Set(res)]
  }

  getPart(): number {
    return this._part
  }

  getSize(x: V): number {
    return this._rank.get(this.find(x)) || 0
  }

  toString(): string {
    return [...this.getGroups()].map(([root, member]) => `${root}: ${member}`).join('\n')
  }
}

if (require.main === module) {
  const uf = new UnionFindArray(10)
  console.log(uf.toString())
  uf.union(0, 1)
  uf.union(1, 2)
  console.log(uf.toString())

  const uf2 = new UnionFindMap()
  console.log(uf2.toString())
  uf2.union('a', 'b')
  uf2.union('b', 'c')
  console.log(uf2.toString())
}

export { UnionFindArray, UnionFindMap }
