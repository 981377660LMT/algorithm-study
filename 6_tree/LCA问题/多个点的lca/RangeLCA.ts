import { Tree } from '../../重链剖分/Tree'

/**
 * 区间LCA(多个点的LCA).
 * @param points 顶点数组.
 * @param lca LCA实现.
 * @returns 返回一个查询 {@link points}[start, end) lca 的函数.
 * @link https://github.com/pranjalssh/CP_codes/blob/master/anta/!RangeLCA.cpp
 */
function rangeLca(
  points: ArrayLike<number>,
  lca: (u: number, v: number) => number
): (start: number, end: number) => number {
  let n = 1
  while (n < points.length) n <<= 1
  const seg = new Int32Array(n << 1)
  for (let i = 0; i < points.length; i++) seg[n + i] = points[i]
  for (let i = n - 1; ~i; i--) seg[i] = lca(seg[i << 1], seg[(i << 1) | 1])

  const calLca = (u: number, v: number): number => {
    if (u === -1 || v === -1) return u === -1 ? v : u
    return lca(u, v)
  }

  const query = (start: number, end: number): number => {
    let res = -1
    for (; start && start + (start & -start) <= end; start += start & -start) {
      res = calLca(res, seg[~~((n + start) / (start & -start))])
    }
    for (; start < end; end -= end & -end) {
      res = calLca(res, seg[~~((n + end) / (end & -end)) - 1])
    }
    return res
  }

  return query
}

export { rangeLca }

if (require.main === module) {
  const tree = new Tree(8)
  tree.addEdge(0, 1)
  tree.addEdge(0, 4)
  tree.addEdge(1, 2)
  tree.addEdge(1, 3)
  tree.addEdge(3, 5)
  tree.addEdge(3, 6)
  tree.addEdge(5, 7)
  tree.build(0)

  const points = Array.from({ length: 8 }, (_, i) => i)
  const query = rangeLca(points, (u, v) => tree.lca(u, v))
  console.log(query(5, 8))
}
