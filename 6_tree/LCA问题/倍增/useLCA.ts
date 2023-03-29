/* eslint-disable no-shadow */
/* eslint-disable no-param-reassign */
/* eslint-disable no-inner-declarations */

function useLCA(n: number, tree: number[][], root = 0) {
  const depth = new Int32Array(n)
  const parent = new Int32Array(n)

  const _bitLen = 32 - Math.clz32(n)
  const _dp = Array(_bitLen).fill(0)
  for (let i = 0; i < _bitLen; i++) {
    _dp[i] = new Int32Array(n).fill(-1)
  }

  _dfs(root, -1, 0)
  _initDp()

  /**
   *  O(logn) 查询LCA
   */
  function queryLCA(root1: number, root2: number): number {
    if (depth[root1] < depth[root2]) [root1, root2] = [root2, root1]

    root1 = upToDepth(root1, depth[root2])
    if (root1 === root2) return root1

    for (let bit = _bitLen - 1; ~bit; bit--) {
      if (_dp[bit][root1] !== _dp[bit][root2]) {
        root1 = _dp[bit][root1]
        root2 = _dp[bit][root2]
      }
    }

    return _dp[0][root1]
  }

  /**
   * O(logn) 查询两点距离
   */
  function queryDist(root1: number, root2: number): number {
    return depth[root1] + depth[root2] - 2 * depth[queryLCA(root1, root2)]
  }

  /**
   * O(logn) 查询树节点root的第k个祖先,如果不存在返回-1
   * @param k k >= 1
   */
  function queryKthAncestor(root: number, k: number): number {
    let bit = 0
    while (k > 0) {
      if (k & 1) {
        root = _dp[root][bit]
        if (root === -1) return -1
      }
      bit++
      k >>>= 1 // 注意:会被转为uint32
    }

    return root
  }

  /**
   * 从`root`开始向上跳至指定深度`toDepth`,返回跳至的节点
   */
  function upToDepth(root: number, toDepth: number): number {
    if (toDepth >= depth[root]) return root
    for (let i = _bitLen; ~i; i--) {
      if ((depth[root] - toDepth) & (1 << i)) {
        root = _dp[i][root]
      }
    }
    return root
  }

  /**
   * 从start节点跳到target节点,跳过step个节点(0-indexed)
   * 返回跳到的节点,如果不存在这样的节点,返回-1
   */
  function jump(start: number, target: number, step: number): number {
    const lca = queryLCA(start, target)
    const dep1 = depth[start]
    const dep2 = depth[target]
    const deplca = depth[lca]
    const dist = dep1 + dep2 - 2 * deplca
    if (step > dist) return -1
    if (step <= dep1 - deplca) return queryKthAncestor(start, step)
    return queryKthAncestor(target, dist - step)
  }

  return {
    depth,
    parent,
    queryLCA,
    queryDist,
    queryKthAncestor,
    upToDepth,
    jump
  }

  function _dfs(cur: number, pre: number, dep: number): void {
    depth[cur] = dep
    parent[cur] = pre
    for (let i = 0; i < tree[cur].length; i++) {
      const next = tree[cur][i]
      if (next !== pre) {
        _dfs(next, cur, dep + 1)
      }
    }
  }

  function _initDp(): void {
    for (let j = 0; j < n; j++) _dp[0][j] = parent[j]
    for (let i = 0; i + 1 < _bitLen; i++) {
      for (let j = 0; j < n; j++) {
        if (_dp[i][j] === -1) _dp[i + 1][j] = -1
        else _dp[i + 1][j] = _dp[i][_dp[i][j]] // 2^i*2^i === 2^(i+1)
      }
    }
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
