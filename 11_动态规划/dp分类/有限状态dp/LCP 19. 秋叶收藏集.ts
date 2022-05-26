/**
 * @param {string} leaves
 * @return {number}
 * 小扣想要将收藏集中树叶的排列调整成「红、黄、红」三部分，每部分树叶数量可以不相等，但均需大于等于 1
 * @summary
 * 注意(leaves[i - 1] === 'r' ? 0 : 1)要加括号
 */
function minimumOperations(leaves: string): number {
  // 最终需要红黄红
  // 目标是「红、黄、红」的话实际上只能通过「红1」「红1、黄」「红1、黄、红2」三种状态转移而来
  let dp = [leaves[0] === 'r' ? 0 : 1, Infinity, Infinity]

  for (let i = 1; i < leaves.length; i++) {
    const ndp = [Infinity, Infinity, Infinity]
    ndp[0] = (leaves[i] === 'r' ? 0 : 1) + dp[0]
    ndp[1] = (leaves[i] === 'y' ? 0 : 1) + Math.min(dp[0], dp[1])
    ndp[2] = (leaves[i] === 'r' ? 0 : 1) + Math.min(dp[1], dp[2])
    dp = ndp
  }

  return dp[2]
}

console.log(minimumOperations('yry'))
export {}
