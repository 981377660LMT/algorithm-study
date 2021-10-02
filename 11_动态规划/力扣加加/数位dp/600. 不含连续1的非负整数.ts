/**
 * 找出小于或等于 n 的非负整数中，其二进制表示不包含 连续的1 的个数。
 * @param n  1 <= n <= 10**9 即树高不超过31
 * 此题极其经典
 * 01 字典树 + 动态规划
 * @description
 * 在 01 字典树的路径上，不能存在连续两个 1 求根节点到叶子节点这样的路径数
 */
function findIntegers(n: number): number {
  /**
   * @description
   * 高为i的满二叉树的路径数量 最下面一层层数最大
   */
  const dp = Array<number>(32).fill(0)
  dp[0] = 1 // 根节点是第0层
  dp[1] = 2 // 00(0) 01(1)
  for (let i = 2; i < 32; i++) {
    // 左子树加右子树的左分支
    dp[i] = dp[i - 1] + dp[i - 2]
  }

  /**
   * 非满二叉树的路径数量
   * @summary
   * 如果该子树有右子树，说明其左子树是满二叉树 先加上满二叉树的数量
   */
  let res = 0
  let pre = 0
  for (let i = 31; ~i; i--) {
    const level = 1 << i
    const hasRight = level & n
    if (hasRight) {
      res += dp[i]
      if (pre === 1) {
        res-- // 减去到根节点0的路径
        break
      }
      pre = 1
    } else {
      pre = 0
    }
  }

  return res + 1 // 到根节点0的路径
}

console.log(findIntegers(5))
// 下面是带有相应二进制表示的非负整数<= 5：
// 0 : 0
// 1 : 1
// 2 : 10
// 3 : 11
// 4 : 100
// 5 : 101
// 其中，只有整数3违反规则（有两个连续的1），其他5个满足规则。
