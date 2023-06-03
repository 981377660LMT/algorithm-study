import assert from 'assert'

/* eslint-disable no-inner-declarations */
import { offLineLca } from '../../../../6_tree/LCA问题/tarjan离线/OfflineLCA'

/**
 * 树上莫队.
 */
class MoOnTree {
  private readonly _tree: number[][]
  private readonly _root: number
  private readonly _queries: [u: number, v: number][] = []

  constructor(tree: number[][], root = 0) {
    this._tree = tree
    this._root = root
  }

  /**
   * 添加从顶点u到顶点v的查询.
   */
  addQuery(u: number, v: number): void {
    this._queries.push([u, v])
  }

  /**
   * 处理每个查询.
   * @param add 将数据添加到窗口.
   * @param remove 将数据从窗口移除.
   * @param query 查询窗口内的数据.
   */
  run(
    add: (rootId: number) => void,
    remove: (rootId: number) => void,
    query: (qid: number) => void
  ): void {
    const n = this._tree.length
    const vs = new Uint32Array(2 * n)
    let vsPtr = 0
    const tin = new Uint32Array(n)
    const tout = new Uint32Array(n)
    const tree = this._tree
    const queries = this._queries
    initTime(this._root, -1)

    const lca = offLineLca(tree, queries, this._root)
    const blockSize = Math.ceil((2 * n) / Math.sqrt(queries.length))
    const qs: [lb: number, l: number, r: number, lca: number, qid: number][] = Array(queries.length)
    queries.forEach(([v, w], i) => {
      if (tin[v] > tin[w]) {
        v ^= w
        w ^= v
        v ^= w
      }
      const lca_ = lca[i]
      if (lca_ !== v) {
        qs[i] = [(tout[v] / blockSize) | 0, tout[v], tin[w] + 1, lca_, i]
      } else {
        qs[i] = [(tin[v] / blockSize) | 0, tin[v], tin[w] + 1, -1, i]
      }
    })

    qs.sort((a, b) => {
      const alb = a[0]
      const blb = b[0]
      if (alb !== blb) return alb - blb
      return alb & 1 ? b[2] - a[2] : a[2] - b[2]
    })

    const flip = new Uint8Array(n)
    let l = 0
    let r = 0
    for (let i = 0; i < qs.length; i++) {
      const [_, ql, qr, lca_, qid] = qs[i]
      while (r < qr) f(vs[r++])
      while (l < ql) f(vs[l++])
      while (l > ql) f(vs[--l])
      while (r > qr) f(vs[--r])
      if (lca_ >= 0) f(lca_)
      query(qid)
      if (lca_ >= 0) f(lca_)
    }

    function initTime(cur: number, pre: number): void {
      tin[cur] = vsPtr
      vs[vsPtr++] = cur
      const tos = tree[cur]
      for (let i = 0; i < tos.length; ++i) {
        const to = tos[i]
        if (to !== pre) initTime(to, cur)
      }
      tout[cur] = vsPtr
      vs[vsPtr++] = cur
    }

    function f(u: number): void {
      flip[u] ^= 1
      if (flip[u]) add(u)
      else remove(u)
    }
  }
}

export { MoOnTree }

if (require.main === module) {
  //   8 2
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
