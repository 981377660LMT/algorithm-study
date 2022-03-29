/**
 * @param {number[]} stones
 * @return {number}
 * 两人不断取前缀和 每次跨度至少+1 取到最后一个前缀和时结束游戏
 * Alice 的目标是 最大化 分数差，Bob 的目标是 最小化 分数差。
 * 返回Alice 和 Bob 的 分数之差 。
 * @summary 我们必须注意到的是，取的所有分数要放回去这个操作，本质是不改变前缀和的。
 */
var stoneGameVIII = function (stones: number[]): number {
  const len = stones.length
  const pre = Array<number>(len + 1).fill(0)
  for (let i = 1; i <= len; i++) {
    pre[i] = pre[i - 1] + stones[i - 1]
  }

  const memo = new Map<number, number>()

  // 起始位置 >=index 时 自己减对方的最大差值
  const dfs = (index: number): number => {
    if (index >= len - 1) return pre[len]
    if (memo.has(index)) return memo.get(index)!
    // 不选index位置:转化为dfs(index+1)
    // 选index:pre[index+1]是自己的分 dfs(index+1) 是对方的分
    const res = Math.max(dfs(index + 1), pre[index + 1] - dfs(index + 1))

    memo.set(index, res)
    return res
  }
  return dfs(1)
}

console.log(stoneGameVIII([-1, 2, -3, 4, -5]))
