/**
 * @param {number} n
 * @return {number}
 * 最初记事本上只有一个字符 'A' 。你每次可以对这个记事本进行两种操作：
 * 你需要使用最少的操作次数，在记事本上输出 恰好 n 个 'A' 。返回能够打印出 n 个 'A' 的最少操作次数。
 * https://leetcode-cn.com/problems/2-keys-keyboard/solution/cong-di-gui-dao-su-shu-fen-jie-by-fuxuemingzhu/
 * @summary 就是让我们求 n 能拆成哪些素因子
 */
var minSteps = function (n: number): number {
  let res = 0
  let prime = 2
  while (n > 1) {
    while (n % prime === 0) {
      res += prime
      n /= prime
    }
    prime++
  }

  return res
}

console.log(minSteps(3))
// 解释：
// 最初, 只有一个字符 'A'。
// 第 1 步, 使用 Copy All 操作。
// 第 2 步, 使用 Paste 操作来获得 'AA'。
// 第 3 步, 使用 Paste 操作来获得 'AAA'。
// 36 = 18 * 2，题目所求的最优结果是 18 + 2 = 20 么
// 因为如果把18 拆开 36 = 3 * 6 * 2，此时复制粘贴的个数只需要 3 + 6 + 2 = 11
// 但这仍然不是最优结果，36 = 3 * 2 * 3 * 2，此时复制粘贴的个数只需要 3 + 2 + 3 + 2 = 10 次。此时已经是最优了。
