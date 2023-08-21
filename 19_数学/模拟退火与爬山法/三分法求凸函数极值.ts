/* eslint-disable no-inner-declarations */

const INF = 2e15

/**
 * 三分法求`严格凸函数fun`在`[lower,upper]`间的最小值.
 */
function minimize(fun: (pos: number) => number, left: number, right: number): number {
  let res = INF
  while (right - left >= 3) {
    const diff = Math.floor((right - left) / 3)
    const pos1 = left + diff
    const pos2 = right - diff
    if (fun(pos1) > fun(pos2)) {
      left = pos1
    } else {
      right = pos2
    }
  }

  while (left <= right) {
    const cand = fun(left)
    res = cand < res ? cand : res
    left++
  }

  return res
}

/**
 * 三分法求`严格凸函数fun`在`[lower,upper]`间的最大值.
 */
function maximize(fun: (pos: number) => number, left: number, right: number): number {
  let res = -INF
  while (right - left >= 3) {
    const diff = Math.floor((right - left) / 3)
    const pos1 = left + diff
    const pos2 = right - diff
    if (fun(pos1) < fun(pos2)) {
      left = pos1
    } else {
      right = pos2
    }
  }

  while (left <= right) {
    const cand = fun(left)
    res = cand > res ? cand : res
    left++
  }

  return res
}

export { minimize, maximize }

if (require.main === module) {
  // 2819. 购买巧克力后的最小相对损失
  // https://leetcode.cn/problems/minimum-relative-loss-after-buying-chocolates/
  function minimumRelativeLosses(prices: number[], queries: number[][]): number[] {
    prices.sort((a, b) => a - b)
    const preSum = Array(prices.length + 1).fill(0)
    const sufSum = Array(prices.length + 1).fill(0)
    for (let i = 1; i <= prices.length; i++) {
      preSum[i] = preSum[i - 1] + prices[i - 1]
      sufSum[i] = sufSum[i - 1] + prices[prices.length - i]
    }

    const res = Array(queries.length).fill(0)
    queries.forEach(([threshold, count], qi) => {
      res[qi] = minimize(
        preLen => {
          const sufLen = count - preLen
          return preSum[preLen] + 2 * threshold * sufLen - sufSum[sufLen]
        },
        0,
        count
      )
    })

    return res
  }
}
