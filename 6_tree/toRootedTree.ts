/**
 * 无根树转有根树.
 */
function toRootedTree(tree: number[][], root = 0): number[][] {
  const n = tree.length
  const rootedTree: number[][] = Array(n)
  for (let i = 0; i < n; i++) rootedTree[i] = []
  const visited = new Uint8Array(n)
  visited[root] = 1
  const queue = new Uint32Array(n)
  let head = 0
  let tail = 0
  queue[tail++] = root
  while (head < tail) {
    const cur = queue[head++]
    const nexts = tree[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i]
      if (!visited[next]) {
        visited[next] = 1
        queue[tail++] = next
        rootedTree[cur].push(next)
      }
    }
  }
  return rootedTree
}

export { toRootedTree }

if (require.main === module) {
  const tree = [[1, 2], [0, 3, 4], [0, 5, 6], [1], [1], [2], [2]]
  console.log(toRootedTree(tree, 0))
}
