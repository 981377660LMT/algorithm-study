/* eslint-disable no-shadow */

// 我对线段树的理解(递归版)
type E = unknown
type Id = unknown
declare const e: <E>() => E
declare const id: <Id>() => Id
declare const op: <E>(data1: E, data2: E) => E
declare const mapping: <E, Id>(data: E, lazy: Id) => E
declare const composition: <Id>(lazy1: Id, lazy2: Id) => Id

class LazySegmentTree {
  private readonly _n: number
  private readonly _data: E[]
  private readonly _lazy: Id[]
  // !别的一些信息 。。。

  constructor(leaves: ArrayLike<unknown>) {
    this._n = leaves.length
    const log = 32 - Math.clz32(this._n - 1)
    const size = 1 << log
    // !初始化data和lazy数组(可用TypedArray优化) 然后建树
    this._data = Array(2 * size).fill(e()) // monoid
    this._lazy = Array(size).fill(id()) // monoid
    this._build(1, 1, this._n, leaves)
  }

  query(left: number, right: number): E {
    return this._query(1, left, right, 1, this._n)
  }

  update(left: number, right: number, lazy: Id): void {
    this._update(1, left, right, 1, this._n, lazy)
  }

  queryAll(): E {
    return this._data[1]
  }

  private _build(root: number, left: number, right: number, leaves: ArrayLike<unknown>): void {
    if (left === right) {
      // !初始化叶子结点data信息
      // this._data[root] = leaves[left - 1]
      return
    }

    const mid = Math.floor((left + right) / 2)
    this._build(root << 1, left, mid, leaves)
    this._build((root << 1) | 1, mid + 1, right, leaves)
    this._pushUp(root, left, right)
  }

  private _query(root: number, L: number, R: number, l: number, r: number): E {
    if (L <= l && r <= R) {
      return this._data[root]
    }

    this._pushDown(root, l, r)
    const mid = Math.floor((l + r) / 2)
    let res = e()
    if (L <= mid) res = op(res, this._query(root << 1, L, R, l, mid))
    if (mid < R) res = op(res, this._query((root << 1) | 1, L, R, mid + 1, r))
    return res
  }

  private _update(root: number, L: number, R: number, l: number, r: number, lazy: Id): void {
    if (L <= l && r <= R) {
      this._propagate(root, l, r, lazy)
      return
    }

    this._pushDown(root, l, r)
    const mid = Math.floor((l + r) / 2)
    if (L <= mid) this._update(root << 1, L, R, l, mid, lazy)
    if (mid < R) this._update((root << 1) | 1, L, R, mid + 1, r, lazy)
    this._pushUp(root, l, r)
  }

  private _pushUp(root: number, left: number, right: number): void {
    // !op操作更新root结点的data信息
    this._data[root] = op(this._data[root << 1], this._data[(root << 1) | 1])
  }

  private _pushDown(root: number, left: number, right: number): void {
    // !传播lazy信息(可以判断根的lazy不为monoid时才传播,传播后将根的lazy置为monoid)
    if (this._lazy[root] !== id()) {
      const mid = Math.floor((left + right) / 2)
      this._propagate(root << 1, left, mid, this._lazy[root])
      this._propagate((root << 1) | 1, mid + 1, right, this._lazy[root])
      this._lazy[root] = id()
    }
  }

  private _propagate(root: number, left: number, right: number, lazy: Id) {
    // !mapping + composition 来更新子节点data和lazy信息
    this._data[root] = mapping(this._data[root], lazy)
    // !判断是否为叶子结点
    if (root < this._lazy.length) {
      this._lazy[root] = composition(this._lazy[root], lazy)
    }
  }
}

export {}
