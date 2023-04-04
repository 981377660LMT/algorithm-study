/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

// 区间并查集 RangeUnionFind/UnionFindRange
// !只使用了路径压缩,每次操作O(logn)

/**
 * 使用路径压缩的区间并查集, 每次操作O(logn).
 */
class UnionFindRange {
  part: number
  private readonly _n: number
  private readonly _parent: Uint32Array
  private readonly _rank: Uint32Array

  constructor(n: number) {
    this.part = n
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
    if (rootX === rootY) {
      return false
    }
    this._parent[rootX] = rootY
    this._rank[rootY] += this._rank[rootX]
    this.part -= 1
    if (f) {
      f(rootY, rootX)
    }
    return true
  }

  /**
   * 合并[left,right]区间, 返回合并次数.
   */
  unionRange(left: number, right: number, f?: (big: number, small: number) => void): number {
    if (left > right) {
      return 0
    }
    const leftRoot = this.find(left)
    let rightRoot = this.find(right)
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
}

/**
 * 利用并查集寻找区间的某个位置左侧/右侧第一个未被访问过的位置.
 * 初始时,所有位置都未被访问过.
 */
class Finder {
  readonly left: Uint32Array // 每个组的左边界,包含端点
  readonly right: Uint32Array // 每个组的右边界,不包含端点
  private readonly _data: Int32Array
  private readonly _n: number // 0 和 n + 1 为哨兵, 实际使用[1,n]

  constructor(n: number) {
    this._n = n
    n += 2
    const data = new Int32Array(n)
    const left = new Uint32Array(n)
    const right = new Uint32Array(n)
    for (let i = 0; i < n; i++) {
      data[i] = -1
      left[i] = i
      right[i] = i + 1
    }
    this._data = data
    this.left = left
    this.right = right
  }

  /**
   * 找到x左侧第一个未被访问过的位置(包含x).
   * 如果不存在, 返回 null.
   */
  prev(x: number): number | null {
    const res = this.left[this.find(x + 1)]
    return res > 0 ? res - 1 : null
  }

  /**
   * 找到x右侧第一个未被访问过的位置(包含x).
   * 如果不存在, 返回 null.
   */
  next(x: number): number | null {
    const res = this.right[this.find(x)]
    return res < this._n + 1 ? res - 1 : null
  }

  /**
   * 删除[start, end)区间内的元素.
   * 0<=start<=end<=n.
   */
  erase(start: number, end: number): number {
    if (start >= end) {
      return
    }
    let count = 0
    while (true) {
      const m = this.right[this.find(start)]
      if (m > end) {
        break
      }
      this.union(start, m)
      count++
    }
    return count
  }

  union(x: number, y: number): boolean {
    let rootX = this.find(x)
    let rootY = this.find(y)
    if (rootX === rootY) {
      return false
    }

    if (this._data[rootX] > this._data[rootY]) {
      rootX ^= rootY
      rootY ^= rootX
      rootX ^= rootY
    }
    this._data[rootX] += this._data[rootY]
    this._data[rootY] = rootX
    this.left[rootX] = Math.min(this.left[rootX], this.left[rootY])
    this.right[rootX] = Math.max(this.right[rootX], this.right[rootY])
    return true
  }

  find(x: number): number {
    if (this._data[x] < 0) {
      return x
    }
    this._data[x] = this.find(this._data[x])
    return this._data[x]
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

export { UnionFindRange, Finder }
