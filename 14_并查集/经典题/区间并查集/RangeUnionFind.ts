/* eslint-disable no-constant-condition */
/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

// 区间并查集 RangeUnionFind/UnionFindRange
// !只使用了路径压缩,每次操作O(logn)

/**
 * 使用路径压缩的区间并查集, 每次操作O(logn).
 */
class UnionFindRange {
  private _part: number
  private readonly _n: number
  private readonly _parent: Uint32Array
  private readonly _rank: Uint32Array

  constructor(n: number) {
    this._part = n
    this._n = n
    this._parent = new Uint32Array(n)
    this._rank = new Uint32Array(n)
    for (let i = 0; i < n; i++) {
      this._parent[i] = i
      this._rank[i] = 1
    }
  }

  find(x: number): number {
    while (x !== this._parent[x]) {
      this._parent[x] = this._parent[this._parent[x]]
      x = this._parent[x]
    }
    return x
  }

  /**
   * union 后, 大的编号所在的组的指向小的编号所在的组.
   */
  union(x: number, y: number, f?: (big: number, small: number) => void): boolean {
    if (x < y) {
      x ^= y
      y ^= x
      x ^= y
    }
    const rootX = this.find(x)
    const rootY = this.find(y)
    if (rootX === rootY) return false
    this._parent[rootX] = rootY
    this._rank[rootY] += this._rank[rootX]
    this._part -= 1
    if (f) f(rootY, rootX)
    return true
  }

  /**
   * 合并[left,right]`闭`区间, 返回合并次数.
   */
  unionRange(start: number, end: number, f?: (big: number, small: number) => void): number {
    if (start >= end) return 0
    const leftRoot = this.find(start)
    let rightRoot = this.find(end)
    let unionCount = 0
    while (rightRoot !== leftRoot) {
      unionCount += 1
      this.union(rightRoot, rightRoot - 1, f)
      rightRoot = this.find(rightRoot - 1)
    }
    return unionCount
  }

  isConnected(x: number, y: number): boolean {
    return this.find(x) === this.find(y)
  }

  getSize(x: number): number {
    return this._rank[this.find(x)]
  }

  getGroups(): Map<number, number[]> {
    const group = new Map<number, number[]>()
    for (let i = 0; i < this._n; i++) {
      const root = this.find(i)
      if (!group.has(root)) {
        group.set(root, [])
      }
      group.get(root)!.push(i)
    }
    return group
  }

  get part(): number {
    return this._part
  }
}

/**
 * 维护每个分组左右边界的区间并查集.
 * 按秩合并.
 */
class UnionFindRange2 {
  readonly groupStart: Uint32Array // 每个组的左边界,包含端点
  readonly groupEnd: Uint32Array // 每个组的右边界,不包含端点
  private _part: number
  private readonly _n: number
  private readonly _data: Int32Array

  constructor(n: number) {
    const data = new Int32Array(n)
    const start = new Uint32Array(n)
    const end = new Uint32Array(n)
    for (let i = 0; i < n; i++) {
      data[i] = -1
      start[i] = i
      end[i] = i + 1
    }
    this.groupStart = start
    this.groupEnd = end
    this._part = n
    this._n = n
    this._data = data
  }

  /**
   * 合并`[start,end)`闭区间, 返回合并次数.
   * 0 <= left <= end < n.
   */
  unionRange(start: number, end: number, f?: (big: number, small: number) => void): number {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return 0
    let count = 0
    while (true) {
      const next = this.groupEnd[this.find(start)]
      if (next >= end) break
      this.union(start, next, f)
      count++
    }
    return count
  }

  union(x: number, y: number, f?: (big: number, small: number) => void): boolean {
    let rootX = this.find(x)
    let rootY = this.find(y)
    if (rootX === rootY) return false
    if (this._data[rootX] > this._data[rootY]) {
      rootX ^= rootY
      rootY ^= rootX
      rootX ^= rootY
    }
    this._data[rootX] += this._data[rootY]
    this._data[rootY] = rootX
    this.groupStart[rootX] = Math.min(this.groupStart[rootX], this.groupStart[rootY])
    this.groupEnd[rootX] = Math.max(this.groupEnd[rootX], this.groupEnd[rootY])
    this._part--
    f && f(rootX, rootY)
    return true
  }

  find(x: number): number {
    if (this._data[x] < 0) return x
    this._data[x] = this.find(this._data[x])
    return this._data[x]
  }

  isConnected(x: number, y: number): boolean {
    return this.find(x) === this.find(y)
  }

  /**
   * 每个点所在分组的左右边界[start,end).
   */
  getRange(x: number): [start: number, end: number] {
    const root = this.find(x)
    return [this.groupStart[root], this.groupEnd[root]]
  }

  getSize(x: number): number {
    return -this._data[this.find(x)]
  }

  getGroups(): Map<number, number[]> {
    const group = new Map<number, number[]>()
    for (let i = 0; i < this._n; i++) {
      const root = this.find(i)
      if (!group.has(root)) {
        group.set(root, [])
      }
      group.get(root)!.push(i)
    }
    return group
  }

  get part(): number {
    return this._part
  }
}

if (require.main === module) {
  // https://leetcode.cn/problems/amount-of-new-area-painted-each-day/
  // 2158. 每天绘制新区域的数量
  function amountPainted(paint: number[][]): number[] {
    const uf = new UnionFindRange(5e4 + 10)
    return paint.map(([left, right]) => uf.unionRange(left, right))
  }
}

export { UnionFindRange, UnionFindRange2 }
