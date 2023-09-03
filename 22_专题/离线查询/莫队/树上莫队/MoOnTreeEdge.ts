/* eslint-disable prefer-destructuring */
/* eslint-disable no-inner-declarations */

import { SortedListRangeBlock } from '../../根号分治/值域分块/SortedListRangeBlock'

/**
 * 维护边权的树上莫队.
 */
class MoOnTreeEdge {
  private readonly _tree: [next: number, eid: number][][]
  private readonly _root: number
  private readonly _queries: [u: number, v: number][] = []

  constructor(tree: [next: number, eid: number][][], root = 0) {
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
   * @param add 将边添加到窗口.
   * @param remove 将边从窗口移除.
   * @param query 查询窗口内的数据.
   */
  run(
    add: (edgeId: number) => void,
    remove: (edgeId: number) => void,
    query: (qid: number) => void
  ): void {
    if (this._queries.length === 0) return

    const n = this._tree.length
    let dfn = 0
    const ins = new Uint32Array(n)
    const outs = new Uint32Array(n)
    const tree = this._tree
    const queries = this._queries
    const dfnToEdge = new Int32Array(2 * n)
    _dfs(this._root, -1)

    const lca = offLineLca(tree, queries, this._root)
    const blockSize = Math.ceil((2 * n) / Math.sqrt(queries.length))
    const qs: [bid: number, left: number, right: number, qid: number][] = Array(queries.length)
    queries.forEach(([v, w], qi) => {
      if (ins[v] > ins[w]) {
        v ^= w
        w ^= v
        v ^= w
      }
      const lca_ = lca[qi]
      if (lca_ !== v) {
        qs[qi] = [(outs[v] / blockSize) | 0, outs[v], ins[w], qi]
      } else {
        qs[qi] = [(ins[v] / blockSize) | 0, ins[v], ins[w], qi]
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
      const [_, ql, qr, qid] = qs[i]
      while (right < qr) _toggle(dfnToEdge[right++])
      while (left < ql) _toggle(dfnToEdge[left++])
      while (left > ql) _toggle(dfnToEdge[--left])
      while (right > qr) _toggle(dfnToEdge[--right])
      query(qid)
    }

    function _dfs(cur: number, pre: number): void {
      ins[cur] = dfn
      const nexts = tree[cur]
      for (let i = 0; i < nexts.length; ++i) {
        const [to, eid] = nexts[i]
        if (to !== pre) {
          dfnToEdge[dfn++] = eid
          _dfs(to, cur)
          dfnToEdge[dfn++] = eid
        }
      }
      outs[cur] = dfn
    }

    function _toggle(u: number): void {
      flip[u] ^= 1
      if (flip[u]) add(u)
      else remove(u)
    }
  }
}

/**
 * 离线求LCA.
 */
function offLineLca(
  tree: [next: number, eid: number][][],
  queries: [start: number, end: number][] | number[][],
  root = 0
): number[] {
  const n = tree.length
  const data = new Int32Array(n)
  const stack = new Uint32Array(n)
  const mark = new Int32Array(n)
  const ptr = new Uint32Array(n)
  let top = 0
  stack[top] = root
  const res = Array(queries.length)
  for (let i = 0; i < queries.length; ++i) res[i] = -1
  const queue: [next: number, ei: number][][] = Array(n)
  for (let i = 0; i < n; ++i) {
    queue[i] = []
    mark[i] = -1
    data[i] = -1
    ptr[i] = tree[i].length
  }
  queries.forEach(([start, end], i) => {
    queue[start].push([end, i])
    queue[end].push([start, i])
  })

  while (top !== -1) {
    const u = stack[top]
    const nexts = tree[u]
    if (mark[u] === -1) {
      mark[u] = u
    } else {
      union(u, nexts[ptr[u]][0])
      mark[find(u)] = u
    }

    if (!run(u)) {
      queue[u].forEach(([v, i]) => {
        if (mark[v] !== -1 && res[i] === -1) {
          res[i] = mark[find(v)]
        }
      })
      --top
    }
  }

  return res

  function union(key1: number, key2: number): boolean {
    let root1 = find(key1)
    let root2 = find(key2)
    if (root1 === root2) return false
    if (data[root1] > data[root2]) {
      root1 ^= root2
      root2 ^= root1
      root1 ^= root2
    }
    data[root1] += data[root2]
    data[root2] = root1
    return true
  }

  function find(key: number): number {
    if (data[key] < 0) return key
    data[key] = find(data[key])
    return data[key]
  }

  function run(u: number): boolean {
    const nexts = tree[u]
    while (ptr[u]) {
      const v = nexts[--ptr[u]][0]
      if (mark[v] === -1) {
        stack[++top] = v
        return true
      }
    }
    return false
  }
}

export { MoOnTreeEdge }

if (require.main === module) {
  // https://leetcode.cn/problems/minimum-edge-weight-equilibrium-queries-in-a-tree/solutions/2424966/typescript-yu-zhi-yu-wu-guan-de-shu-shan-k42t/
  function minOperationsQueries(n: number, edges: number[][], queries: number[][]): number[] {
    const tree: [next: number, eid: number][][] = Array(n)
    for (let i = 0; i < n; i++) tree[i] = []
    edges.forEach(([u, v], ei) => {
      tree[u].push([v, ei])
      tree[v].push([u, ei])
    })

    const mo = new MoOnTreeEdge(tree, 0)
    queries.forEach(([u, v]) => mo.addQuery(u, v))

    const res: number[] = Array(queries.length).fill(0)
    const weightCounter = new Map<number, number>()
    const sl = new SortedListRangeBlock(n - 1)
    let edgeCount = 0

    const add = (eid: number): void => {
      const weight = edges[eid][2]
      const preCount = weightCounter.get(weight) || 0
      weightCounter.set(weight, preCount + 1)
      sl.discard(preCount)
      sl.add(preCount + 1)
      edgeCount++
    }

    const remove = (eid: number): void => {
      const weight = edges[eid][2]
      const preCount = weightCounter.get(weight)!
      weightCounter.set(weight, preCount - 1)
      sl.discard(preCount)
      if (preCount - 1) sl.add(preCount - 1)
      edgeCount--
    }

    const query = (qid: number): void => {
      if (!sl.length) return
      res[qid] = edgeCount - sl.max!
    }

    mo.run(add, remove, query)
    return res
  }

  console.log(
    minOperationsQueries(
      8,
      [
        [1, 2, 6],
        [1, 3, 4],
        [2, 4, 6],
        [2, 5, 3],
        [3, 6, 6],
        [3, 0, 8],
        [7, 0, 2]
      ],
      [
        [4, 6],
        [0, 4],
        [6, 5],
        [7, 4]
      ]
    )
  )

  console.log(minOperationsQueries(1, [], [[0, 0]]))
}
