/**
 * 获取子树信息.
 * height[i] 表示以 i 为根的子树的高度(距离最远的叶子节点的距离).
 */
function getSubtreeInfo(
  tree: number[][],
  root: number
): {
  parent: Int32Array
  depth: Int32Array
  subsize: Int32Array
  height: Int32Array
} {
  const n = tree.length
  const parent = new Int32Array(n)
  const depth = new Int32Array(n)
  const subsize = new Int32Array(n)
  const height = new Int32Array(n)
  const topological = new Int32Array(n)
  topological[0] = root
  parent[root] = root
  depth[root] = 0
  let r = 1
  for (let l = 0; l < r; l++) {
    const cur = topological[l]
    const nexts = tree[cur]
    for (let j = 0; j < nexts.length; j++) {
      const next = nexts[j]
      if (next !== parent[cur]) {
        topological[r] = next
        r++
        parent[next] = cur
        depth[next] = depth[cur] + 1
      }
    }
  }

  for (r--; r >= 0; r--) {
    const cur = topological[r]
    subsize[cur] = 1
    height[cur] = 0
    const nexts = tree[cur]
    for (let j = 0; j < nexts.length; j++) {
      const next = nexts[j]
      if (next !== parent[cur]) {
        subsize[cur] += subsize[next]
        height[cur] = Math.max(height[cur], height[next] + 1)
      }
    }
  }

  parent[root] = -1

  return { parent, depth, subsize, height }
}

export { getSubtreeInfo }

if (require.main === module) {
  // tree := [][]int32{{1, 2}, {0}, {0, 3, 4}, {2}, {2, 5}, {4}}
  const tree = [[1, 2], [0], [0, 3, 4], [2], [2, 5], [4]]
  const root = 0
  const { parent, depth, subsize, height } = getSubtreeInfo(tree, root)
  console.log(parent, depth, subsize, height)
}
