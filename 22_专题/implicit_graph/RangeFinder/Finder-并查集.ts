/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */
// 寻找前驱后继/区间删除

import { onlineBfs } from '../OnlineBfs-在线bfs'

/**
 * 利用并查集寻找区间的某个位置左侧/右侧第一个未被访问过的位置.
 * 初始时,所有位置都未被访问过.
 */
class Finder {
  private readonly _n: number
  private readonly _lParent: Uint32Array
  private readonly _rParent: Uint32Array

  constructor(n: number) {
    this._n = n
    const lp = new Uint32Array(n + 2)
    const rp = new Uint32Array(n + 2)
    for (let i = 0; i < lp.length; i++) {
      lp[i] = i
      rp[i] = i
    }
    this._lParent = lp
    this._rParent = rp
  }

  /**
   * 找到x左侧第一个未被访问过的位置(包含x).
   * 如果不存在, 返回 null.
   */
  prev(x: number): number | null {
    const res = this._lFind(x + 1)
    return res > 0 ? res - 1 : null
  }

  /**
   * 找到x右侧第一个未被访问过的位置(包含x).
   * 如果不存在, 返回 null.
   */
  next(x: number): number | null {
    const res = this._rFind(x + 1)
    return res < this._n + 1 ? res - 1 : null
  }

  /**
   * 删除[start, end)区间内的元素.
   * 0<=start<=end<=n.
   */
  erase(start: number, end: number): void {
    if (start >= end) {
      return
    }
    let leftRoot = this._rFind(start + 1)
    let rightRoot = this._rFind(end + 1)
    while (rightRoot !== leftRoot) {
      this._rUnion(leftRoot, leftRoot + 1)
      leftRoot = this._rFind(leftRoot + 1)
    }
    leftRoot = this._lFind(start)
    rightRoot = this._lFind(end)
    while (rightRoot !== leftRoot) {
      this._lUnion(rightRoot, rightRoot - 1)
      rightRoot = this._lFind(rightRoot - 1)
    }
  }

  private _lUnion(x: number, y: number): void {
    if (x < y) {
      ;[x, y] = [y, x]
    }
    const rootX = this._lFind(x)
    const rootY = this._lFind(y)
    if (rootX === rootY) {
      return
    }
    this._lParent[rootX] = rootY
  }

  private _rUnion(x: number, y: number): void {
    if (x > y) {
      ;[x, y] = [y, x]
    }
    const rootX = this._rFind(x)
    const rootY = this._rFind(y)
    if (rootX === rootY) {
      return
    }
    this._rParent[rootX] = rootY
  }

  private _lFind(x: number): number {
    while (x !== this._lParent[x]) {
      this._lParent[x] = this._lParent[this._lParent[x]]
      x = this._lParent[x]
    }
    return x
  }

  private _rFind(x: number): number {
    while (x !== this._rParent[x]) {
      this._rParent[x] = this._rParent[this._rParent[x]]
      x = this._rParent[x]
    }
    return x
  }
}

if (require.main === module) {
  // https://leetcode.cn/problems/minimum-reverse-operations/submissions/

  function minReverseOperations(n: number, p: number, banned: number[], k: number): number[] {
    const finder = [new Finder(n), new Finder(n)]
    for (let i = 0; i < n; i++) {
      finder[(i & 1) ^ 1].erase(i, i + 1)
    }
    banned.forEach(i => {
      finder[i & 1].erase(i, i + 1)
    })

    const getRange = (i: number): [number, number] => [
      Math.max(i - k + 1, k - i - 1),
      Math.min(i + k - 1, 2 * n - k - i - 1)
    ]

    const setUsed = (u: number): void => {
      finder[u & 1].erase(u, u + 1)
    }

    const findUnused = (u: number): number | null => {
      const [left, right] = getRange(u)
      const next = finder[(u + k + 1) & 1].next(left)
      if (next !== null && left <= next && next <= right) {
        return next
      }
      return null
    }

    const dist = onlineBfs(n, p, setUsed, findUnused)
    return dist.map(d => (d === INF ? -1 : d))
  }
}

export { Finder }
