/**
 * @param {string} leaves
 * @return {number}
 * 小扣想要将收藏集中树叶的排列调整成「红、黄、红」三部分
 * @summary
 * leetcode股票问题 dp[i][j] i表示日期 j表示状态
 * 注意(leaves[i - 1] === 'r' ? 0 : 1)要加括号
 */
var minimumOperations = function (leaves: string): number {
  // 最终需要红黄红
  // 目标是「红、黄、红」的话实际上只能通过「红1」「红1、黄」「红1、黄、红2」三种状态转移而来
  const dp = Array.from({ length: leaves.length }, () => [Infinity, Infinity, Infinity])
  dp[0] = [leaves[0] === 'r' ? 0 : 1, Infinity, Infinity]
  dp[1][2] = Infinity

  for (let i = 1; i < leaves.length; i++) {
    dp[i][0] = dp[i - 1][0] + (leaves[i] === 'r' ? 0 : 1)
    dp[i][1] = Math.min(dp[i - 1][0], dp[i - 1][1]) + (leaves[i] === 'y' ? 0 : 1)
    dp[i][2] = Math.min(dp[i - 1][1], dp[i - 1][2]) + (leaves[i] === 'r' ? 0 : 1)
  }

  return dp[leaves.length - 1][2]
}
// class Solution:
//     def minimumOperations(self, leaves: str) -> int:
//         # 最终需要红黄红
//         # 维持3种状态，分别为截止目前全部红色，在截止目前全红的基础上变为黄色，以及变为黄色的基础上全部红色
//         n = len(leaves)
//         dp = [[float('inf')]*3 for _ in range(n)]
//         for i in range(n):
//             if i ==0:
//                 dp[i][0] = 0 if leaves[i]=='r' else 1
//                 continue
//             dp[i][0] = min(dp[i][0],dp[i-1][0] + (leaves[i]!='r'))
//             dp[i][1] = min(dp[i-1][1]+(leaves[i]!='y'),dp[i-1][0] +(leaves[i]!='y') )
//             dp[i][2] = min(dp[i-1][2]+(leaves[i]!='r'),dp[i-1][1]+(leaves[i]!='r'))
//         return dp[n-1][2]

// 作者：dabien1
// 链接：https://leetcode-cn.com/problems/UlBDOe/solution/zhe-dao-ti-ru-guo-qia-zhu-de-jian-yi-sou-suo-yi-xi/
// 来源：力扣（LeetCode）
// 著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。
// console.log(minimumOperations('rrryyyrryyyrr'))
console.log(minimumOperations('yry'))
