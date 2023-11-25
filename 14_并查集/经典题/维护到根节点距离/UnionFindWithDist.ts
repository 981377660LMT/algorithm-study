/* eslint-disable max-len */
/* eslint-disable no-param-reassign */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

// API:
//  union(x,y,dist) : p(x) = p(y) + dist. 如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
//  find(x) : 返回x所在组的根节点.
//  dist(x,y) : 返回x到y的距离.
//  distToRoot(x) : 返回x到所在组根节点的距离.

/**
 * 维护到根节点距离的并查集.
 * 用于维护环的权值，树上的距离等.
 * !Array实现.
 */
class UnionFindWithDistArray<V> {
  private _part: number
  private readonly _data: Int32Array
  private readonly _potential: V[]
  private readonly _e: () => V
  private readonly _op: (a: V, b: V) => V
  private readonly _inv: (a: V) => V

  constructor(
    n: number,
    monoid: {
      e: () => V
      op: (a: V, b: V) => V
      inv: (a: V) => V
    } & ThisType<void>
  ) {
    this._part = n
    this._data = new Int32Array(n).fill(-1)
    this._potential = Array(n)
    for (let i = 0; i < n; i++) this._potential[i] = monoid.e()
    this._e = monoid.e
    this._op = monoid.op
    this._inv = monoid.inv
  }

  /**
   * p(child) = p(parent) + dist.
   * 如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
   */
  union(child: number, parent: number, dist: V, cb?: (big: number, small: number) => void): boolean {
    dist = this._op(dist, this._op(this.distToRoot(parent), this._inv(this.distToRoot(child))))
    child = this.find(child)
    parent = this.find(parent)
    if (child === parent) {
      return dist === this._e()
    }
    if (this._data[child] < this._data[parent]) {
      child ^= parent
      parent ^= child
      child ^= parent
      dist = this._inv(dist)
    }
    this._data[parent] += this._data[child]
    this._data[child] = parent
    this._potential[child] = dist
    this._part--
    cb && cb(parent, child)
    return true
  }

  find(x: number): number {
    if (this._data[x] < 0) {
      return x
    }
    const root = this.find(this._data[x])
    this._potential[x] = this._op(this._potential[x], this._potential[this._data[x]])
    this._data[x] = root
    return root
  }

  /**
   * 返回x到y的距离`f(x) - f(y)`.
   */
  dist(x: number, y: number): V {
    return this._op(this.distToRoot(x), this._inv(this.distToRoot(y)))
  }

  /**
   * 返回x到所在组根节点的距离`f(x) - f(find(x))`.
   */
  distToRoot(x: number): V {
    this.find(x)
    return this._potential[x]
  }

  isConnected(x: number, y: number): boolean {
    return this.find(x) === this.find(y)
  }

  getSize(x: number): number {
    return -this._data[this.find(x)]
  }

  getGroups(): Map<number, number[]> {
    const res: Map<number, number[]> = new Map()
    for (let i = 0; i < this._data.length; i++) {
      const root = this.find(i)
      if (!res.has(root)) res.set(root, [])
      res.get(root)!.push(i)
    }
    return res
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

/**
 * 维护到根节点距离的并查集.
 * 用于维护环的权值，树上的距离等.
 * !Map实现.
 */
class UnionFindWithDistMap<V> {
  private _part = 0
  private readonly _data = new Map<number, number>()
  private readonly _potential = new Map<number, V>()
  private readonly _e: () => V
  private readonly _op: (a: V, b: V) => V
  private readonly _inv: (a: V) => V

  constructor(
    monoid: {
      e: () => V
      op: (a: V, b: V) => V
      inv: (a: V) => V
    } & ThisType<void>
  ) {
    this._e = monoid.e
    this._op = monoid.op
    this._inv = monoid.inv
  }

  /**
   * p(child) = p(parent) + dist.
   * 如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
   */
  union(child: number, parent: number, dist: V, cb?: (big: number, small: number) => void): boolean {
    dist = this._op(dist, this._op(this.distToRoot(parent), this._inv(this.distToRoot(child))))
    child = this.find(child)
    parent = this.find(parent)
    if (child === parent) {
      return dist === this._e()
    }
    if (this._data.get(child)! < this._data.get(parent)!) {
      child ^= parent
      parent ^= child
      child ^= parent
      dist = this._inv(dist)
    }
    this._data.set(parent, this._data.get(child)! + this._data.get(parent)!)
    this._data.set(child, parent)
    this._potential.set(child, dist)
    this._part--
    cb && cb(parent, child)
    return true
  }

  find(x: number): number {
    if (!this._data.has(x)) {
      this.add(x)
      return x
    }
    if (this._data.get(x)! < 0) {
      return x
    }
    const root = this.find(this._data.get(x)!)
    this._potential.set(x, this._op(this._potential.get(x)!, this._potential.get(this._data.get(x)!)!))
    this._data.set(x, root)
    return root
  }

  /**
   * 返回x到y的距离`f(x) - f(y)`.
   */
  dist(x: number, y: number): V {
    return this._op(this.distToRoot(x), this._inv(this.distToRoot(y)))
  }

  /**
   * 返回x到所在组根节点的距离`f(x) - f(find(x))`.
   */
  distToRoot(x: number): V {
    this.find(x)
    return this._potential.get(x)!
  }

  isConnected(x: number, y: number): boolean {
    return this.find(x) === this.find(y)
  }

  getSize(x: number): number {
    return -this._data.get(this.find(x))!
  }

  getGroups(): Map<number, number[]> {
    const res: Map<number, number[]> = new Map()
    this._data.forEach((_, key) => {
      const root = this.find(key)
      if (!res.has(root)) res.set(root, [])
      res.get(root)!.push(key)
    })
    return res
  }

  add(x: number): this {
    if (!this._data.has(x)) {
      this._data.set(x, -1)
      this._potential.set(x, this._e())
      this._part++
    }
    return this
  }

  has(x: number): boolean {
    return this._data.has(x)
  }

  get part(): number {
    return this._part
  }
}

export { UnionFindWithDistArray, UnionFindWithDistMap }

if (require.main === module) {
  const uf = new UnionFindWithDistArray<number>(10, { e: () => 0, op: (a, b) => a + b, inv: a => -a })
  uf.union(0, 1, 1)
  uf.union(1, 3, 2)
  uf.union(4, 5, 3)
  console.log(uf.dist(0, 1))
}
