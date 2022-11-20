/* eslint-disable no-shadow */

// 我对线段树的理解(递归版)
type Data = unknown
type Lazy = unknown
declare const e: <Data>() => Data
declare const id: <Lazy>() => Lazy
declare const op: <Data>(data1: Data, data2: Data) => Data
declare const mapping: <Data, Lazy>(data: Data, lazy: Lazy) => Data
declare const composition: <Lazy>(lazy1: Lazy, lazy2: Lazy) => Lazy

class LazySegmentTree {
  private readonly _n: number
  private readonly _data: Data[]
  private readonly _lazy: Lazy[]
  // !别的一些信息 。。。

  constructor(leaves: ArrayLike<unknown>) {
    this._n = leaves.length
    const cap = 1 << (32 - Math.clz32(this._n - 1) + 1)
    // !初始化data和lazy数组(可用TypedArray优化) 然后建树
    this._data = Array(cap).fill(e()) // monoid
    this._lazy = Array(cap).fill(id()) // monoid
    this._build(1, 1, this._n, leaves)
  }

  query(left: number, right: number): Data {
    return this._query(1, left, right, 1, this._n)
  }

  update(left: number, right: number, lazy: Lazy): void {
    this._update(1, left, right, 1, this._n, lazy)
  }

  queryAll(): Data {
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

  private _query(root: number, L: number, R: number, l: number, r: number): Data {
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

  private _update(root: number, L: number, R: number, l: number, r: number, lazy: Lazy): void {
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

  private _propagate(root: number, left: number, right: number, lazy: Lazy) {
    // !mapping + composition 来更新子节点data和lazy信息
    this._data[root] = mapping(this._data[root], lazy)
    this._lazy[root] = composition(this._lazy[root], lazy)
  }
}

export {}
