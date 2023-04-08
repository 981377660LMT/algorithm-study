/**
 * 需要高速化 dp[pos][使用次数] 的 dp 时,
 * 如果dp(k+1)-dp(k)<=dp(k)-dp(k-1) ，则可以使用 wqs 二分.
 * !问题转化为`每使用一次操作罚款 penalty 元,求最大分数`.
 * 对penalty 二分搜索，转化为 dp[pos]一个维度的dp.
 *
 * @param k 最大操作次数
 * @param getDp dp: func(penalty int) [2]int: 每使用一次操作罚款 penalty 元, 返回 [子问题dp的`最大值`, `最大的`操作使用次数]
 *
 * @see
 * https://leetcode-cn.com/problems/best-time-to-buy-and-sell-stock-iv/solution/yi-chong-ji-yu-wqs-er-fen-de-you-xiu-zuo-x36r/
 */
function aliensDp(
  k: number,
  getDp: (penalty: number) => [best: number, maxOperation: number]
): number {
  let left = 1
  let right = 2e15
  let penalty = 0

  while (left <= right) {
    const mid = (left + right) >>> 1
    const cand = getDp(mid)
    if (cand[1] >= k) {
      penalty = mid
      left = mid + 1
    } else {
      right = mid - 1
    }
  }

  const res = getDp(penalty)
  return res[0] + penalty * k
}

export { aliensDp }
