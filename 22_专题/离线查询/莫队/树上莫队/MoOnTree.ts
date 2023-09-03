/* eslint-disable no-inner-declarations */

import assert from 'assert'

import { offLineLca } from '../../../../6_tree/LCA问题/tarjan离线/OfflineLCA'
import { SortedListRangeBlock } from '../../根号分治/值域分块/SortedListRangeBlock'

/**
 * 维护点权的树上莫队.
 */
class MoOnTree {
  private readonly _tree: number[][]
  private readonly _root: number
  private readonly _queries: [from: number, to: number][] = []

  constructor(tree: number[][], root = 0) {
    this._tree = tree
    this._root = root
  }

  /**
   * 添加从顶点from到顶点to的查询.
   */
  addQuery(from: number, to: number): void {
    this._queries.push([from, to])
  }

  /**
   * 处理每个查询.
   * @param add 将顶点添加到窗口.
   * @param remove 将顶点从窗口移除.
   * @param query 查询窗口内的数据.
   */
  run(
    add: (rootId: number) => void,
    remove: (rootId: number) => void,
    query: (qid: number) => void
  ): void {
    if (this._queries.length === 0) return

    const n = this._tree.length
    let dfn = 0
    const dfnToNode = new Uint32Array(2 * n)
    const ins = new Uint32Array(n)
    const outs = new Uint32Array(n)
    const tree = this._tree
    const queries = this._queries
    _dfs(this._root, -1)

    const lca = offLineLca(tree, queries, this._root)
    const blockSize = Math.ceil((2 * n) / Math.sqrt(queries.length))
    const qs: [lb: number, l: number, r: number, lca: number, qid: number][] = Array(queries.length)
    queries.forEach(([v, w], i) => {
      if (ins[v] > ins[w]) {
        v ^= w
        w ^= v
        v ^= w
      }
      const lca_ = lca[i]
      if (lca_ !== v) {
        qs[i] = [(outs[v] / blockSize) | 0, outs[v], ins[w] + 1, lca_, i]
      } else {
        qs[i] = [(ins[v] / blockSize) | 0, ins[v], ins[w] + 1, -1, i]
      }
    })

    qs.sort((a, b) => {
      const id1 = a[0]
      const id2 = b[0]
      if (id1 !== id2) return id1 - id2
      return id1 & 1 ? b[2] - a[2] : a[2] - b[2]
    })

    const flip = new Uint8Array(n)
    let left = 0
    let right = 0
    for (let i = 0; i < qs.length; i++) {
      const [_, ql, qr, lca_, qid] = qs[i]
      while (right < qr) _toggle(dfnToNode[right++])
      while (left < ql) _toggle(dfnToNode[left++])
      while (left > ql) _toggle(dfnToNode[--left])
      while (right > qr) _toggle(dfnToNode[--right])
      if (lca_ >= 0) _toggle(lca_)
      query(qid)
      if (lca_ >= 0) _toggle(lca_)
    }

    function _dfs(cur: number, pre: number): void {
      ins[cur] = dfn
      dfnToNode[dfn++] = cur
      const tos = tree[cur]
      for (let i = 0; i < tos.length; ++i) {
        const to = tos[i]
        if (to !== pre) _dfs(to, cur)
      }
      outs[cur] = dfn
      dfnToNode[dfn++] = cur
    }

    function _toggle(u: number): void {
      flip[u] ^= 1
      if (flip[u]) add(u)
      else remove(u)
    }
  }
}

export { MoOnTree }

if (require.main === module) {
  路径上的顶点种类数()
  function 路径上的顶点种类数(): void {
    // 8 2
    // 105 2 9 3 8 5 7 7
    // 1 2
    // 1 3
    // 1 4
    // 3 5
    // 3 6
    // 3 7
    // 4 8
    // 2 5
    // 3 8

    const n = 8
    const values = [105, 2, 9, 3, 8, 5, 7, 7]
    const edges = [
      [0, 1],
      [0, 2],
      [0, 3],
      [2, 4],
      [2, 5],
      [2, 6],
      [3, 7]
    ]
    const adjList: number[][] = Array(n)
    for (let i = 0; i < n; i++) adjList[i] = []
    edges.forEach(([u, v]) => {
      adjList[u].push(v)
      adjList[v].push(u)
    })
    const queries = [
      [1, 4],
      [2, 7],
      [6, 7],
      [2, 5]
    ]

    assert.deepStrictEqual(solve(adjList, values, queries), [4, 4, 4, 2])
    console.log('passed')

    // 路径上的顶点种类数
    function solve(
      tree: number[][],
      values: number[],
      queries: [from: number, to: number][] | number[][]
    ): number[] {
      // 离散化
      const pool = new Map<unknown, number>()
      function id(o: unknown): number {
        const res = pool.get(o)
        if (res !== void 0) return res
        const cur = pool.size
        pool.set(o, cur)
        return cur
      }

      const n = tree.length
      const q = queries.length
      const newValue = new Uint32Array(n)
      for (let i = 0; i < n; ++i) newValue[i] = id(values[i])

      const M = new MoOnTree(tree)
      queries.forEach(([u, v]) => {
        M.addQuery(u, v)
      })

      const res: number[] = Array(q)
      const counter = new Uint32Array(pool.size)
      let cur = 0
      M.run(add, remove, query)
      return res

      function add(u: number): void {
        const x = newValue[u]
        if (!counter[x]) cur++
        counter[x]++
      }
      function remove(u: number): void {
        const x = newValue[u]
        if (counter[x] === 1) cur--
        counter[x]--
      }
      function query(qid: number): void {
        res[qid] = cur
      }
    }
  }
}
