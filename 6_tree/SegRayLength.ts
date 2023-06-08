/* eslint-disable prefer-destructuring */

/**
 * 用于处理树中最长射线/线段的问题.
 * @see {@link https://github.com/old-yan/CP-template/blob/main/TREE/SegRayLengthHelper_vector.h}
 */
class SegRayLength {
  /**
   * 以i为端点的线段长度前三大的值.
   */
  readonly ray: [max1: number, max2: number, max3: number][]

  /**
   * 不经过i的线段长度前两大的值.
   */
  readonly seg: [max1: number, max2: number][]

  /**
   * 以i为端点的向下的线段长度最大值.
   */
  readonly downRay: number[]

  /**
   * i下方不越过i的线段长度最大值.
   */
  readonly downSeg: number[]

  /**
   * 以i为端点的向上的线段长度最大值.
   */
  readonly upRay: number[]

  /**
   * i上方不越过i的线段长度最大值.
   */
  readonly upSeg: number[]

  private readonly _tree: [to: number, weight: number][][]
  private readonly _depthWeighted: number[]
  private readonly _lid: Uint32Array
  private readonly _top: Uint32Array
  private readonly _parent: Int32Array
  private readonly _heavySon: Int32Array
  private _dfn = 0

  constructor(n: number) {
    const ray = Array(n)
    const seg = Array(n)
    const downRay = Array(n)
    const downSeg = Array(n)
    const upRay = Array(n)
    const upSeg = Array(n)
    const tree = Array(n)
    const depthWeighted = Array(n)
    const lid = new Uint32Array(n)
    const top = new Uint32Array(n)
    const parent = new Int32Array(n)
    const heavySon = new Int32Array(n)
    for (let i = 0; i < n; ++i) {
      ray[i] = [0, 0, 0]
      seg[i] = [0, 0]
      downRay[i] = 0
      downSeg[i] = 0
      upRay[i] = 0
      upSeg[i] = 0
      tree[i] = []
      depthWeighted[i] = 0
      parent[i] = -1
    }

    this.ray = ray
    this.seg = seg
    this.downRay = downRay
    this.downSeg = downSeg
    this.upRay = upRay
    this.upSeg = upSeg
    this._tree = tree
    this._depthWeighted = depthWeighted
    this._lid = lid
    this._top = top
    this._parent = parent
    this._heavySon = heavySon
  }

  addEdge(u: number, v: number, w = 1): void {
    this._tree[u].push([v, w])
    this._tree[v].push([u, w])
  }

  addDirectedEdge(u: number, v: number, w = 1): void {
    this._tree[u].push([v, w])
  }

  build(root = 0): void {
    this._dfs1(root, -1, 0)
    this._dfs2(root, 0, 0)
    this._dfs3(root, root)
  }

  /**
   * 查询最长射线ray,最长线段seg.
   * @param u 树节点.
   * @param ignoreRoot 屏蔽掉以ignoreRoot为根的子树.需要保证ignoreRoot在u的子树中.
   * @returns 树中剩余部分从u出发的最长射线,树中剩余部分的最长线段.
   */
  queryMaxRayAndSeg(u: number, ignoreRoot = -1): [maxRay: number, maxSeg: number] {
    if (u === ignoreRoot) return this._maxRaySeg(u, -1, -1)
    return this._maxRaySeg(
      u,
      this.downRay[ignoreRoot] + this._weightedDist(u, ignoreRoot),
      this.downSeg[ignoreRoot]
    )
  }

  private _dfs1(cur: number, pre: number, dist: number): number {
    let subSize = 1
    let heavySize = 0
    let heavySon = -1
    const treeCur = this._tree[cur]
    for (let i = 0; i < treeCur.length; ++i) {
      const [next, weight] = treeCur[i]
      if (next !== pre) {
        const nextSize = this._dfs1(next, cur, dist + weight)
        subSize += nextSize
        if (nextSize > heavySize) {
          heavySize = nextSize
          heavySon = next
        }
        this._addDownRay(cur, this.downRay[next] + weight)
        this._addDownSeg(cur, this.downSeg[next])
      }
    }
    const [len1, len2] = this.ray[cur]
    const cand = len1 + len2
    if (cand > this.downSeg[cur]) this.downSeg[cur] = cand
    this._depthWeighted[cur] = dist
    this._parent[cur] = pre
    this._heavySon[cur] = heavySon
    return subSize
  }

  private _dfs2(cur: number, curRay: number, curSeg: number): void {
    this._setUpRay(cur, curRay)
    this._setUpSeg(cur, curSeg)
    const treeCur = this._tree[cur]
    for (let i = 0; i < treeCur.length; ++i) {
      const [next, weight] = treeCur[i]
      if (next !== this._parent[cur]) {
        let [ray, seg] = this._maxRaySeg(cur, this.downRay[next] + weight, this.downSeg[next])
        this._addSeg(next, seg)
        ray += weight
        if (ray > seg) seg = ray
        this._dfs2(next, ray, seg)
      }
    }
  }

  private _dfs3(cur: number, top: number): void {
    this._top[cur] = top
    this._lid[cur] = this._dfn++
    const heavySon = this._heavySon[cur]
    if (~heavySon) {
      this._dfs3(heavySon, top)
      this._tree[cur].forEach(([next]) => {
        if (next !== heavySon && next !== this._parent[cur]) {
          this._dfs3(next, next)
        }
      })
    }
  }

  private _addRay(i: number, ray: number): void {
    const rayI = this.ray[i]
    if (ray > rayI[0]) {
      rayI[2] = rayI[1]
      rayI[1] = rayI[0]
      rayI[0] = ray
    } else if (ray > rayI[1]) {
      rayI[2] = rayI[1]
      rayI[1] = ray
    } else if (ray > rayI[2]) {
      rayI[2] = ray
    }
  }

  private _addSeg(i: number, seg: number): void {
    const segI = this.seg[i]
    if (seg > segI[0]) {
      segI[1] = segI[0]
      segI[0] = seg
    } else if (seg > segI[1]) {
      segI[1] = seg
    }
  }

  private _addDownRay(i: number, ray: number): void {
    if (ray > this.downRay[i]) this.downRay[i] = ray
    this._addRay(i, ray)
  }

  private _addDownSeg(i: number, seg: number): void {
    if (seg > this.downSeg[i]) this.downSeg[i] = seg
    this._addSeg(i, seg)
  }

  private _setUpRay(i: number, ray: number): void {
    this.upRay[i] = ray
    this._addRay(i, ray)
  }

  private _setUpSeg(i: number, seg: number): void {
    this.upSeg[i] = seg
  }

  /**
   * 查询树中某部分的最长射线和线段.
   * @param u 树中某点.
   * @param ignoreRay 屏蔽掉的部分提供的最长射线.
   * @param ignorSeg 屏蔽掉的部分提供的最长线段.
   * @returns 从u出发的最长射线和树中剩余部分的最长线段.
   */
  private _maxRaySeg(u: number, ignoreRay: number, ignorSeg: number): [number, number] {
    const [r0, r1, r2] = this.ray[u]
    const [s0, s1] = this.seg[u]
    const maxRay = ignoreRay === r0 ? r1 : r0
    let maxSeg = ignorSeg === s0 ? s1 : s0
    let twoRay = r0 + r1
    if (ignoreRay === r0) twoRay = r1 + r2
    else if (ignoreRay === r1) twoRay = r0 + r2
    if (maxSeg < twoRay) maxSeg = twoRay
    return [maxRay, maxSeg]
  }

  private _lca(u: number, v: number): number {
    while (true) {
      if (this._lid[u] > this._lid[v]) {
        u ^= v
        v ^= u
        u ^= v
      }
      if (this._top[u] === this._top[v]) return u
      v = this._parent[this._top[v]]
    }
  }

  private _weightedDist(u: number, v: number): number {
    return (
      this._depthWeighted[u] + this._depthWeighted[v] - 2 * this._depthWeighted[this._lca(u, v)]
    )
  }
}

export { SegRayLength }

if (require.main === module) {
  const SR = new SegRayLength(5)
  SR.addEdge(2, 0, 1)
  SR.addEdge(1, 3, 1)
  SR.addEdge(4, 0, 1)
  SR.addEdge(3, 0, 1)
  SR.build(3)

  console.log(SR.ray)
  console.log(SR.seg)
  console.log(SR.downRay)
  console.log(SR.downSeg)
  console.log(SR.upRay)
  console.log(SR.upSeg)

  console.log(SR.queryMaxRayAndSeg(0, 2)) // (2, 3)

  // eslint-disable-next-line no-inner-declarations
  function treeDiameter(edges: number[][]): number {
    const SR = new SegRayLength(edges.length + 1)
    edges.forEach(([u, v]) => SR.addEdge(u, v, 1))
    SR.build()
    return SR.ray.reduce((max, [r0]) => Math.max(max, r0), 0)
  }
}
