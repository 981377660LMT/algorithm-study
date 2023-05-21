/* eslint-disable no-param-reassign */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

// API:
//  union(x,y,dist) : p(x) = p(y) + dist. 如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
//  find(x) : 返回x所在组的根节点.
//  dist(x,y) : 返回x到y的距离.
//  distToRoot(x) : 返回x到所在组根节点的距离.

/**
 * 维护到根节点距离的并查集.
 * !距离为加法.Array实现.
 */
class UnionFindArrayWithDist {
  private _part: number
  private readonly _data: Int32Array
  private readonly _potential: number[]

  constructor(n: number) {
    this._part = n
    this._data = new Int32Array(n).fill(-1)
    this._potential = Array(n).fill(0)
  }

  /**
   * p(x) = p(y) + dist.
   * 如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
   */
  union(x: number, y: number, dist: number, cb?: (big: number, small: number) => void): boolean {
    dist += this.distToRoot(y) - this.distToRoot(x)
    x = this.find(x)
    y = this.find(y)
    if (x === y) {
      return dist === 0
    }
    if (this._data[x] < this._data[y]) {
      x ^= y
      y ^= x
      x ^= y
      dist = -dist
    }
    this._data[y] += this._data[x]
    this._data[x] = y
    this._potential[x] = dist
    this._part--
    cb && cb(y, x)
    return true
  }

  find(x: number): number {
    if (this._data[x] < 0) {
      return x
    }
    const root = this.find(this._data[x])
    this._potential[x] += this._potential[this._data[x]]
    this._data[x] = root
    return root
  }

  /**
   * 返回x到y的距离`f(x) - f(y)`.
   */
  dist(x: number, y: number): number {
    return this.distToRoot(x) - this.distToRoot(y)
  }

  /**
   * 返回x到所在组根节点的距离`f(x) - f(find(x))`.
   */
  distToRoot(x: number): number {
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
      if (!res.has(root)) {
        res.set(root, [])
      }
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
 * !距离为加法.Map实现.
 */
class UnionFindMapWithDist {
  private _part = 0
  private readonly _data = new Map<number, number>()
  private readonly _potential = new Map<number, number>()

  /**
   * p(x) = p(y) + dist.
   * 如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
   */
  union(x: number, y: number, dist: number, cb?: (big: number, small: number) => void): boolean {
    dist += this.distToRoot(y) - this.distToRoot(x)
    x = this.find(x)
    y = this.find(y)
    if (x === y) {
      return dist === 0
    }
    if (this._data.get(x)! < this._data.get(y)!) {
      x ^= y
      y ^= x
      x ^= y
      dist = -dist
    }
    this._data.set(y, this._data.get(x)! + this._data.get(y)!)
    this._data.set(x, y)
    this._potential.set(x, dist)
    this._part--
    cb && cb(y, x)
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
    this._potential.set(x, this._potential.get(x)! + this._potential.get(this._data.get(x)!)!)
    this._data.set(x, root)
    return root
  }

  /**
   * 返回x到y的距离`f(x) - f(y)`.
   */
  dist(x: number, y: number): number {
    return this.distToRoot(x) - this.distToRoot(y)
  }

  /**
   * 返回x到所在组根节点的距离`f(x) - f(find(x))`.
   */
  distToRoot(x: number): number {
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
      if (!res.has(root)) {
        res.set(root, [])
      }
      res.get(root)!.push(key)
    })
    return res
  }

  add(x: number): this {
    if (!this._data.has(x)) {
      this._data.set(x, -1)
      this._potential.set(x, 0)
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

export { UnionFindArrayWithDist, UnionFindMapWithDist }
