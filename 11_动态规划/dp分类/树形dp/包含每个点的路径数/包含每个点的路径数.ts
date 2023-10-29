/**
 * 求包含每个点的路径数.路径至少有两个点.
 */
function countPath(n: number, tree: ArrayLike<ArrayLike<number>>): number[] {
  const res = Array<number>(n).fill(0)
  const dfs = (cur: number, pre: number): number => {
    let count = 0
    let size = 1
    const nexts = tree[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i]
      if (next !== pre) {
        const subSize = dfs(next, cur)
        count += size * subSize
        size += subSize
      }
    }
    count += size * (n - size)
    res[cur] = count
    return size
  }
  dfs(0, -1)
  return res
}

export { countPath }
