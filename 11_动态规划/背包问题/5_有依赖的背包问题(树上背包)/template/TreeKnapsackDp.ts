type Item = { value: number; weight: number }

function treeKnapsackDpNaive(tree: ArrayLike<ArrayLike<number>>, items: ArrayLike<Item>, maxCapacity: number, root = 0): number {
  const dfs = (cur: number, pre: number): number[] => {
    const { value, weight } = items[cur]
    const dp = Array(maxCapacity + 1).fill(0)
    dp.fill(value, weight, maxCapacity + 1) // 根节点必须选
    const nexts = tree[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i]
      if (next === pre) continue
      const ndp = dfs(next, cur)
      for (let j = maxCapacity; j >= weight; j--) {
        // 类似分组背包，枚举分给子树 to 的容量 w，对应的子树的最大价值为 dt[w]
        // w 不可超过 j-it.weight，否则无法选择根节点
        for (let w = 0; w <= j - weight; w++) {
          dp[j] = Math.max(dp[j], dp[j - w] + ndp[w])
        }
      }
    }
    return dp
  }

  const res = dfs(root, -1)
  return res[maxCapacity]
}

function treeKnapsackDpSquare(tree: ArrayLike<ArrayLike<number>>, items: ArrayLike<Item>, maxCapacity: number, root = 0): number {}

export { treeKnapsackDpSquare, treeKnapsackDpNaive }
