/* eslint-disable no-inner-declarations */

import { WaveletMatrixSegments } from '../../../24_高级数据结构/waveletmatrix/WaveletMatrixForTree'

interface Options {
  tree: Tree

  /**
   * 顶点的值, 或者边的值.(边的编号为两个定点中较深的那个点的编号).
   */
  data: Uint32Array

  /**
   * data是否为顶点的值以及查询的时候是否是顶点权值.
   */
  isVertex: boolean

  /**
   * 如果要支持异或,则需要按照异或的值来决定值域.设为-1时表示不使用异或.
   */
  log?: number

  /**
   * 如果要支持区间和,设为空时表示不使用区间和.
   */
  sumData?: number[]
}

/**
 * 树的路径查询, 维护的量需要满足群的性质，且必须要满足交换律.
 */
class TreeWaveletMatrix {
  private readonly _isVertex: boolean
  private readonly _tree: Tree
  private readonly _wm: WaveletMatrixSegments

  constructor(options: Options) {
    const { tree, data, isVertex, log = -1, sumData = [] } = options

    const n = tree.tree.length
    this._tree = tree
    this._isVertex = isVertex
    const A1 = new Uint32Array(n)
    const S1 = sumData.length ? Array(n).fill(0) : []

    if (isVertex) {
      for (let v = 0; v < n; v++) {
        A1[tree.lid[v]] = data[v]
      }
      if (sumData.length === n) {
        for (let v = 0; v < n; v++) {
          S1[tree.lid[v]] = sumData[v]
        }
      }
    } else {
      if (sumData.length) {
        for (let e = 0; e < n - 1; e++) {
          S1[tree.lid[tree.eidTov(e)]] = sumData[e]
        }
      }
      for (let e = 0; e < n - 1; e++) {
        A1[tree.lid[tree.eidTov(e)]] = data[e]
      }
    }

    this._wm = new WaveletMatrixSegments(A1, log, S1)
  }

  /**
   * s到t的路径上有多少个数在[a,b)之间.
   */
  countPath(s: number, t: number, a: number, b: number, xor = 0): number {
    return this._wm.countRangeSegments(this._getSegments(s, t), a, b, xor)
  }

  /**
   * 子树root中有多少个数在[a,b)之间.
   */
  countSubtree(root: number, a: number, b: number, xor = 0): number {
    const l = this._tree.lid[root]
    const r = this._tree.rid[root]
    const offset = this._isVertex ? 0 : 1
    return this._wm.countRange(l + offset, r, a, b, xor)
  }

  /**
   * s到t的路径上第k小的数以及前k个数的和.
   */
  kthValueAndSumPath(
    s: number,
    t: number,
    k: number,
    xor = 0
  ): [res: number, preSum: number] | [res: null, preSum: number] {
    return this._wm.kthValueAndSumSegments(this._getSegments(s, t), k, xor)
  }

  /**
   * 子树内第k小的数以及前k个数的和.
   */
  kthValueAndSumSubtree(
    root: number,
    k: number,
    xor = 0
  ): [res: number, preSum: number] | [res: null, preSum: number] {
    const l = this._tree.lid[root]
    const r = this._tree.rid[root]
    const offset = this._isVertex ? 0 : 1
    return this._wm.kthValueAndSum(l + offset, r, k, xor)
  }

  /**
   * s到t的路径上第k小的数.
   */
  kthPath(s: number, t: number, k: number, xor = 0): number {
    return this._wm.kthSegments(this._getSegments(s, t), k, xor)
  }

  /**
   * 子树内第k小的数.
   */
  kthSubtree(root: number, k: number, xor = 0): number {
    const l = this._tree.lid[root]
    const r = this._tree.rid[root]
    const offset = this._isVertex ? 0 : 1
    return this._wm.kth(l + offset, r, k, xor)
  }

  /**
   * s到t的路径上的中位数.
   */
  medianPath(upper: boolean, s: number, t: number, xor = 0): number {
    return this._wm.medianSegments(upper, this._getSegments(s, t), xor)
  }

  /**
   * 子树内的中位数.
   */
  medianSubtree(upper: boolean, root: number, xor = 0): number {
    const l = this._tree.lid[root]
    const r = this._tree.rid[root]
    const offset = this._isVertex ? 0 : 1
    return this._wm.median(upper, l + offset, r, xor)
  }

  /**
   * s到t的路径上第k1小到第k2小的数的和.
   */
  sumPath(s: number, t: number, k1: number, k2: number, xor = 0): number {
    return this._wm.sumSegments(this._getSegments(s, t), k1, k2, xor)
  }

  /**
   * 子树内第k1小到第k2小的数的和.
   */
  sumSubtree(root: number, k1: number, k2: number, xor = 0): number {
    const l = this._tree.lid[root]
    const r = this._tree.rid[root]
    const offset = this._isVertex ? 0 : 1
    return this._wm.sum(l + offset, r, k1, k2, xor)
  }

  /**
   * s到t的路径上所有数的和.
   */
  sumAllPath(s: number, t: number): number {
    return this._wm.sumAllSegments(this._getSegments(s, t))
  }

  /**
   * 子树内所有数的和.
   */
  sumAllSubtree(root: number): number {
    const l = this._tree.lid[root]
    const r = this._tree.rid[root]
    const offset = this._isVertex ? 0 : 1
    return this._wm.sumAll(l + offset, r)
  }

  private _getSegments(s: number, t: number): [number, number][] {
    const segments = this._tree.getPathDecomposition(s, t, this._isVertex)
    for (let i = 0; i < segments.length; i++) {
      const seg = segments[i]
      if (seg[0] > seg[1]) {
        ;[seg[0], seg[1]] = [seg[1], seg[0]]
      }
      seg[1]++
    }
    return segments
  }
}

class Tree {
  readonly tree: [next: number, weight: number][][]
  readonly depth: Uint32Array
  readonly parent: Int32Array
  readonly depthWeighted: number[]
  readonly edges: [from: number, to: number, weight: number][] = []
  readonly lid: Uint32Array
  readonly rid: Uint32Array
  private readonly _idToNode: Uint32Array
  private readonly _top: Uint32Array
  private readonly _heavySon: Int32Array
  private _timer = 0

  constructor(n: number) {
    this.tree = Array(n).fill(0)
    for (let i = 0; i < n; i++) {
      this.tree[i] = []
    }
    this.depth = new Uint32Array(n)
    this.parent = new Int32Array(n).fill(-1)
    this.lid = new Uint32Array(n)
    this.rid = new Uint32Array(n)
    this._idToNode = new Uint32Array(n)
    this._top = new Uint32Array(n)
    this.depthWeighted = Array(n).fill(0)
    this._heavySon = new Int32Array(n)
  }

  addEdge(from: number, to: number, weight = 1): void {
    this.tree[from].push([to, weight])
    this.tree[to].push([from, weight])
    this.edges.push([from, to, weight])
  }

  addDirectedEdge(from: number, to: number, weight = 1): void {
    this.tree[from].push([to, weight])
    this.edges.push([from, to, weight])
  }

  build(root = -1): void {
    if (root === -1) {
      for (let i = 0; i < this.tree.length; i++) {
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

  id(root: number): [inId: number, outId: number] {
    return [this.lid[root], this.rid[root]]
  }

  eid(u: number, v: number): number {
    const id1 = this.lid[u]
    const id2 = this.lid[v]
    return id1 > id2 ? id1 : id2
  }

  eidTov(eid: number): number {
    const e = this.edges[eid]
    const u = e[0]
    const v = e[1]
    return this.parent[u] === v ? u : v
  }

  lca(u: number, v: number): number {
    while (true) {
      if (this.lid[u] > this.lid[v]) {
        ;[u, v] = [v, u]
      }
      if (this._top[u] === this._top[v]) {
        return u
      }
      v = this.parent[this._top[v]]
    }
  }

  dist(u: number, v: number, weighted: boolean): number {
    if (weighted) {
      return this.depthWeighted[u] + this.depthWeighted[v] - 2 * this.depthWeighted[this.lca(u, v)]
    }
    return this.depth[u] + this.depth[v] - 2 * this.depth[this.lca(u, v)]
  }

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
    this.tree[root].forEach(([next]) => {
      if (next !== this.parent[root]) {
        res.push(next)
      }
    })
    return res
  }

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

  subtreeSize(v: number, root = -1): number {
    if (root === -1) {
      return this.rid[v] - this.lid[v]
    }
    if (v === root) {
      return this.tree.length
    }
    const x = this.jump(v, root, 1)
    if (this.isInSubtree(v, x)) {
      return this.rid[v] - this.lid[v]
    }
    return this.tree.length - this.rid[x] + this.lid[x]
  }

  isInSubtree(child: number, root: number): boolean {
    return this.lid[root] <= this.lid[child] && this.rid[child] <= this.rid[root]
  }

  getHeavyPath(start: number): number[] {
    const res: number[] = [start]
    let cur = start
    while (~this._heavySon[cur]) {
      cur = this._heavySon[cur]
      res.push(cur)
    }
    return res
  }

  private _build(cur: number, pre: number, dep: number, dist: number): number {
    let subSize = 1
    let heavySon = -1
    let heavySize = 0
    this.tree[cur].forEach(([next, weight]) => {
      if (next !== pre) {
        const nextSize = this._build(next, cur, dep + 1, dist + weight)
        subSize += nextSize
        if (nextSize > heavySize) {
          heavySize = nextSize
          heavySon = next
        }
      }
    })
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
    if (~this._heavySon[cur]) {
      this._markTop(this._heavySon[cur], top)
      this.tree[cur].forEach(([next]) => {
        if (next !== this._heavySon[cur] && next !== this.parent[cur]) {
          this._markTop(next, next)
        }
      })
    }
    this.rid[cur] = this._timer
  }
}

if (require.main === module) {
  const n = 100000
  const q = 100000

  console.time('build')
  const tree = new Tree(n)
  for (let i = 1; i < n; i++) {
    tree.addEdge(Math.floor(Math.random() * i), i)
  }

  tree.build(0)

  const twm = new TreeWaveletMatrix({
    tree,
    data: new Uint32Array(n),
    isVertex: true
  })

  for (let i = 0; i < q; i++) {
    const u = Math.floor(Math.random() * n)
    const v = Math.floor(Math.random() * n)
    twm.kthPath(u, v, 100, 0)
  }
  console.timeEnd('build')
}

export { TreeWaveletMatrix }
