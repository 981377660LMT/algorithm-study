/* eslint-disable no-shadow */
/* eslint-disable no-param-reassign */
/* eslint-disable no-inner-declarations */

function useLCA(n: number, adjList: number[][], root = 0) {
  const depth = new Int32Array(n).fill(-1)
  const parent = new Int32Array(n).fill(-1)

  const _bitLen = Math.floor(Math.log2(n)) + 1
  const _fa = Array.from<unknown, Int32Array>({ length: n }, () => new Int32Array(_bitLen))

  _dfs(root, -1, 0)
  _initDp()

  /**
   *  O(logn) 查询LCA
   */
  function queryLCA(root1: number, root2: number): number {
    if (depth[root1] < depth[root2]) [root1, root2] = [root2, root1]

    for (let bit = _bitLen - 1; ~bit; bit--) {
      if (depth[_fa[root1][bit]] >= depth[root2]) {
        root1 = _fa[root1][bit]
      }
    }

    if (root1 === root2) return root1

    for (let bit = _bitLen - 1; ~bit; bit--) {
      if (_fa[root1][bit] !== _fa[root2][bit]) {
        root1 = _fa[root1][bit]
        root2 = _fa[root2][bit]
      }
    }

    return _fa[root1][0]
  }

  /**
   * O(logn) 查询两点距离
   */
  function queryDist(root1: number, root2: number): number {
    return depth[root1] + depth[root2] - 2 * depth[queryLCA(root1, root2)]
  }

  /**
   * O(logn) 查询树节点root的第k个祖先,如果不存在返回-1
   */
  function queryKthAncestor(root: number, k: number): number {
    let bit = 0
    while (k) {
      if (k & 1) {
        root = _fa[root][bit]
        if (root === -1) return -1
      }

      k >>= 1
      bit++
    }

    return root
  }

  function _dfs(cur: number, pre: number, dep: number): void {
    depth[cur] = dep
    parent[cur] = pre
    for (let i = 0; i < adjList[cur].length; i++) {
      const next = adjList[cur][i]
      if (next === pre) continue
      _dfs(next, cur, dep + 1)
    }
  }

  /**
   * @description O(nlogn) 初始化dp
   */
  function _initDp(): void {
    for (let i = 0; i < n; i++) _fa[i][0] = parent[i]
    for (let bit = 0; bit < _bitLen - 1; bit++) {
      for (let i = 0; i < n; i++) {
        if (_fa[i][bit] === -1) _fa[i][bit + 1] = -1
        else _fa[i][bit + 1] = _fa[_fa[i][bit]][bit]
      }
    }
  }

  return {
    depth,
    parent,
    queryLCA,
    queryDist,
    queryKthAncestor
  }
}

if (require.main === module) {
  function closestNode(n: number, edges: number[][], query: number[][]): number[] {
    const adjList = Array.from<unknown, number[]>({ length: n }, () => [])
    edges.forEach(([u, v]) => {
      adjList[u].push(v)
      adjList[v].push(u)
    })

    const { depth, queryLCA } = useLCA(n, adjList)
    const res = []
    for (const [root1, root2, root3] of query) {
      const cands = [queryLCA(root1, root2), queryLCA(root2, root3), queryLCA(root1, root3)].sort(
        (a, b) => -(depth[a] - depth[b])
      )
      res.push(cands[0])
    }

    return res
  }
}

export { useLCA }
