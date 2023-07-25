// https://ei1333.github.io/library/graph/tree/centroid-decomposition.hpp
// https://ei1333.github.io/library/test/verify/aoj-3139.test.cpp
// https://ei1333.github.io/library/test/verify/yosupo-frequency-table-of-tree-distance.test.cpp
// https://ei1333.github.io/library/test/verify/yukicoder-1002.test.cpp

// 重心互相连接形成的有根树, 可以想象把树拎起来, 重心在树的中心，连接着各个子树的重心...
//            3 (重)
//          / | \
//     (重)1  0  2 (重)
//        / \    |
//       4   5   6

/**
 * 树的重心分解, 返回点分树和点分树的根.
 * @param n 节点数.
 * @param tree `无向树`的邻接表.
 * @returns [centTree: 重心互相连接形成的`有根树`, root: 点分树的根]
 */
function centroidDecomposition(
  n: number,
  tree: [next: number, weight: number][][]
): [centTree: number[][], root: number] {
  const subSize = new Uint32Array(n)
  const removed = new Uint8Array(n)
  const centTree: number[][] = Array(n)
  for (let i = 0; i < n; i++) centTree[i] = []

  const getSize = (cur: number, parent: number): number => {
    subSize[cur] = 1
    const nexts = tree[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i][0]
      if (next !== parent && !removed[next]) {
        subSize[cur] += getSize(next, cur)
      }
    }
    return subSize[cur]
  }

  const getCentroid = (cur: number, parent: number, mid: number): number => {
    const nexts = tree[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i][0]
      if (next === parent || removed[next]) continue
      if (subSize[next] > mid) return getCentroid(next, cur, mid)
    }
    return cur
  }

  const build = (cur: number): number => {
    const centroid = getCentroid(cur, -1, getSize(cur, -1) >>> 1)
    removed[centroid] = 1
    const nexts = tree[centroid]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i][0]
      if (!removed[next]) {
        centTree[centroid].push(build(next))
      }
    }
    removed[centroid] = 0
    return centroid
  }

  const root = build(0)
  return [centTree, root]
}

export { centroidDecomposition }
