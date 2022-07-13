function useLCAManager(n: number, adjMap: Map<number, Set<number>>, root = 0) {
  const depth = new Int32Array(n).fill(-1)
  const parent = new Int32Array(n).fill(-1)

  const _BITLEN = Math.floor(Math.log2(n)) + 1
  const _fa = Array.from<any, Int32Array>({ length: n }, () => new Int32Array(_BITLEN))

  _dfs(root, -1, 0)
  _initDp()

  /**
   * @description O(logn) 查询LCA
   */
  function queryLCA(root1: number, root2: number): number {
    if (depth[root1] < depth[root2]) [root1, root2] = [root2, root1]

    for (let bit = _BITLEN - 1; ~bit; bit--) {
      if (depth[_fa[root1][bit]] >= depth[root2]) {
        root1 = _fa[root1][bit]
      }
    }

    if (root1 === root2) return root1

    for (let bit = _BITLEN - 1; ~bit; bit--) {
      if (_fa[root1][bit] !== _fa[root2][bit]) {
        root1 = _fa[root1][bit]
        root2 = _fa[root2][bit]
      }
    }

    return _fa[root1][0]
  }

  function _dfs(cur: number, pre: number, dep: number): void {
    depth[cur] = dep
    parent[cur] = pre
    for (const next of adjMap.get(cur) ?? []) {
      if (next === pre) continue
      _dfs(next, cur, dep + 1)
    }
  }

  /**
   * @description O(nlogn) 初始化dp
   */
  function _initDp(): void {
    for (let i = 0; i < n; i++) _fa[i][0] = parent[i]
    for (let bit = 0; bit < _BITLEN - 1; bit++) {
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
  }
}

if (require.main === module) {
  function closestNode(n: number, edges: number[][], query: number[][]): number[] {
    const adjMap = new Map<number, Set<number>>()
    for (const [u, v] of edges) {
      adjMap.set(u, (adjMap.get(u) ?? new Set()).add(v))
      adjMap.set(v, (adjMap.get(v) ?? new Set()).add(u))
    }

    const { depth, queryLCA } = useLCAManager(n, adjMap)
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

export { useLCAManager }
