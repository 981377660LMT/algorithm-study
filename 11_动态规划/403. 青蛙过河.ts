/**
 * @param {number[]} stones
 * @return {boolean}
 * 请判定青蛙能否成功过河:即能否在最后一步跳至最后一块石子上
 * 初始条件：第 0 个石头可以跳 1 步
 * 如果青蛙上一步跳跃了 k 个单位，
 * 那么它接下来的跳跃距离只能选择为 k - 1、k 或 k + 1 个单位。 另请注意，青蛙只能向前方（终点的方向）跳跃。
 */
const canCross = function (stones: number[]): boolean {
  if (stones[1] !== 1) {
    return false
  }
  // dp 里是的状态是石头索引和跳到此石头上时跳跃一次的步数 那么dfs的变量就是这两个
  // 即第i个石头走j步
  const dp = Array.from<unknown, boolean[]>({ length: stones.length }, () =>
    Array(stones.length + 1).fill(false)
  )

  dp[0][1] = true
  for (let i = 1; i < stones.length; i++) {
    // 从i前面的石头j跳
    for (let j = 0; j < i; j++) {
      const step = stones[i] - stones[j]
      if (step <= j + 1) {
        dp[i][step] = dp[j][step - 1] || dp[j][step] || dp[j][step + 1]
        // 提前结束循环直接返回结果
        if (i === stones.length - 1 && dp[i][step]) return true
      }
    }
  }

  return false
}

console.log(canCross([0, 1, 3, 5, 6, 8, 12, 17]))
// 输入：stones = [0,1,3,5,6,8,12,17]
// 输出：true
// 解释：青蛙可以成功过河，按照如下方案跳跃：
// 跳 1 个单位到第 2 块石子, 然后跳 2 个单位到第 3 块石子,
// 接着 跳 2 个单位到第 4 块石子, 然后跳 3 个单位到第 6 块石子,
// 跳 4 个单位到第 7 块石子, 最后，跳 5 个单位到第 8 个石子（即最后一块石子）。

// 错过买卖股票最佳时机，倾家荡产，为了生计，
// 做过粉刷房子的工人，给栅栏涂过色，
// 只有零钱兑换一点食物，无奈还是变为打家劫舍的小偷，
// 通过不同路径来到河边，放下背包，看到青蛙，
// 原地架锅起火，回想动态规划的一生！
