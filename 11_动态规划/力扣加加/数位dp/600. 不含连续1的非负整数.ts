/**
 * 找出小于或等于 n 的非负整数中，其二进制表示不包含 连续的1 的个数。
 * @param n  1 <= n <= 10**9 即树高不超过31
 * 此题极其经典
 * 01 字典树 + 动态规划
 * @description
 * 在 01 字典树的路径上，不能存在连续两个 1 求根节点到叶子节点这样的路径数
 * 1 <= n <= 109
 */
function findIntegers(n: number): number {
  const len = calLength(n)
  // 前多少位，最后一位为0还是1
  const dp = Array.from<unknown, number[]>({ length: len + 1 }, () => [0, 0])
  dp[0][0] = 1

  // 左边i大的是高位
  for (let i = 1; i < len + 1; i++) {
    dp[i][0] = dp[i - 1][0] + dp[i - 1][1]
    dp[i][1] = dp[i - 1][0]
  }

  let res = 0
  let preBit = 0 // 左一位是0还是1
  for (let i = len - 1; ~i; i--) {
    if (n & (1 << i)) {
      res += dp[i + 1][0] // 因此如果该位填 0 的话，后面的低位填什么都是满足要求的 例如 0010101=>0000101
      if (preBit === 1) break

      preBit = 1
    } else {
      // 如果当前位 cur = 0，我们只能选 0，并决策下一位。
      preBit = 0
    }

    if (i === 0) {
      res++
    }
  }

  return res

  function calLength(n: number) {
    let res = 0

    while (n) {
      res++
      n >>>= 1
    }

    return res
  }

  return n.toString(2).length
}
if (require.main === module) {
  console.log(findIntegers(5))
}
// 下面是带有相应二进制表示的非负整数<= 5：
// 0 : 0
// 1 : 1
// 2 : 10
// 3 : 11
// 4 : 100
// 5 : 101
// 其中，只有整数3违反规则（有两个连续的1），其他5个满足规则。

// 结论：事实上长度为n的二进制数字中不含连续1的非负整数个数满足斐波那契数列
// dp[n] = dp[n - 2] + dp[n - 1]
// 该递推式表示长度为n的字符串，都可以通过以下两种方式获得：
// 在长度为n-2的字符串左侧加10；
// 在长度为n-1的字符串左侧加0。
export { findIntegers }
