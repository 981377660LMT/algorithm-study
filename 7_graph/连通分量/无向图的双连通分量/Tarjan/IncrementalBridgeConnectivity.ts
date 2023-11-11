// 维护图中桥的个数，支持增加边的操作
// https://ei1333.github.io/library/graph/connected-components/incremental-bridge-connectivity.hpp
// 概要
// 辺の追加クエリのみ存在するとき, 二重辺連結成分を効率的に管理するデータ構造.
//
// 使い方
// IncrementalBridgeConnectivity(sz): sz 頂点で初期化する.
// addEdge(x, y): 頂点 x と y との間に無向辺を追加する.
// countBridge(): 現在の橋の個数を返す.
// find(k): 頂点 k が属する二重辺連結成分(の代表元)を求める.

import { UnionFindArray } from '../../../../14_并查集/UnionFind'

/**
 * @link https://scrapbox.io/data-structures/Incremental_Bridge-Connectivity
 */
class IncrementalBridgeConnectivity {
  private readonly _cc: UnionFindArray
  private readonly _bcc: UnionFindArray
  private readonly _bbf: Uint32Array
  private _bridgeCount = 0

  constructor(n: number) {
    this._cc = new UnionFindArray(n)
    this._bcc = new UnionFindArray(n)
    this._bbf = new Uint32Array(n).fill(n)
  }

  /**
   * 添加一条无向边.
   */
  addEdge(x: number, y: number): void {
    x = this._bcc.find(x)
    y = this._bcc.find(y)
    if (this._cc.find(x) === this._cc.find(y)) {
      const w = this._lca(x, y)
      this._compress(x, w)
      this._compress(y, w)
    } else {
      if (this._cc.getSize(x) > this._cc.getSize(y)) {
        const tmp = x
        x = y
        y = tmp
      }
      this._link(x, y)
      this._cc.union(x, y)
      this._bridgeCount++
    }
  }

  /**
   * 返回 k 所在的双连通分量的代表元.
   */
  find(k: number): number {
    return this._bcc.find(k)
  }

  /**
   * 返回当前的桥的个数.
   */
  countBridge(): number {
    return this._bridgeCount
  }

  private _size(): number {
    return this._bbf.length
  }

  private _parent(x: number): number {
    if (this._bbf[x] === this._size()) return this._size()
    return this._bcc.find(this._bbf[x])
  }

  private _lca(x: number, y: number): number {
    const used = new Set<number>()
    while (true) {
      if (x !== this._size()) {
        if (used.has(x)) return x
        used.add(x)
        x = this._parent(x)
      }
      const tmp = x
      x = y
      y = tmp
    }
  }

  private _compress(x: number, y: number): void {
    while (this._bcc.find(x) !== this._bcc.find(y)) {
      const next = this._parent(x)
      this._bbf[x] = this._bbf[y]
      this._bcc.union(x, y)
      x = next
      this._bridgeCount--
    }
  }

  private _link(x: number, y: number): void {
    let v = x
    let pre = y
    while (v !== this._size()) {
      const next = this._parent(v)
      this._bbf[v] = pre
      pre = v
      v = next
    }
  }
}

export { IncrementalBridgeConnectivity }

if (require.main === module) {
  const n = 5
  const ibc = new IncrementalBridgeConnectivity(n)
  ibc.addEdge(0, 1)
  ibc.addEdge(1, 2)
  ibc.addEdge(2, 3)
  ibc.addEdge(3, 4)
  console.log(ibc.countBridge()) // 4
  ibc.addEdge(0, 4)
  console.log(ibc.countBridge()) // 0
  console.log(ibc.find(0), ibc.find(1), ibc.find(2), ibc.find(3), ibc.find(4)) // 0 0 0 0 0
}
