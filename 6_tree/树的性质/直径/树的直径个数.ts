/**
 * 树的直径个数.
 */
function countDiameter(
  tree: ArrayLike<ArrayLike<number>>,
  start = 0
): [diameter: number, diameterCount: number] {
  let diameter = 0
  let diameterCount = 0
  const dfs = (cur: number, pre: number): { depth: number; count: number } => {
    let maxDepth = 0
    let count = 1
    const nexts = tree[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i]
      if (next === pre) continue
      const { depth, count: nextCount } = dfs(next, cur)
      if (maxDepth + depth > diameter) {
        diameter = maxDepth + depth
        diameterCount = count * nextCount
      } else if (maxDepth + depth === diameter) {
        diameterCount += count * nextCount
      }
      if (depth > maxDepth) {
        maxDepth = depth
        count = nextCount
      } else if (depth === maxDepth) {
        count += nextCount
      }
    }

    return { depth: maxDepth + 1, count }
  }

  dfs(start, -1)
  return [diameter, diameterCount]
}

export { countDiameter }

if (require.main === module) {
  console.log(countDiameter([[1], [0, 2], [1, 3, 4], [2], [2]]))
}
