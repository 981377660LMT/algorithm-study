/* eslint-disable prefer-destructuring */
/* eslint-disable no-inner-declarations */
/* eslint-disable max-len */

const INF = 2e15

/**
 * 斐波那契搜索寻找`[left,right]`中的一个极值点,不要求单峰性质.
 * @param fun 代价函数.
 * @param left 搜索区间左端点(包含).
 * @param right 搜索区间右端点(包含).
 * @param min 是否寻找最小值.
 * @returns 极值点的横坐标x和纵坐标f(x).
 */
function fibonacciSearch(
  fun: (x: number) => number,
  left: number,
  right: number,
  min: boolean
): [x: number, y: number] {
  let a = left
  let b = left + 1
  let c = left + 2
  let d = left + 3
  let step = 0
  while (d < right) {
    b = c
    c = d
    d = b + c - a
    step++
  }

  const at = (i: number): number => {
    if (right < i) return INF
    return min ? fun(i) : -fun(i)
  }

  let ya = at(a)
  let yb = at(b)
  let yc = at(c)
  let yd = at(d)
  for (let i = 0; i < step; i++) {
    if (yb < yc) {
      d = c
      c = b
      b = a + d - c
      yd = yc
      yc = yb
      yb = at(b)
    } else {
      a = b
      b = c
      c = a + d - b
      ya = yb
      yb = yc
      yc = at(c)
    }
  }

  let x = a
  let y = ya
  if (yb < y) {
    x = b
    y = yb
  }
  if (yc < y) {
    x = c
    y = yc
  }
  if (yd < y) {
    x = d
    y = yd
  }

  return min ? [x, y] : [x, -y]
}

/** 斐波那契搜索求`凸函数fun`在`[left,right]`间的`(极小值点,极小值)`. */
const minimize = (fun: (x: number) => number, left: number, right: number) => fibonacciSearch(fun, left, right, true)

/** 斐波那契搜索求`凹函数fun`在`[left,right]`间的`(极大值点,极大值)`. */
const maximize = (fun: (x: number) => number, left: number, right: number) => fibonacciSearch(fun, left, right, false)

export { fibonacciSearch, minimize, maximize }

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

    const res: number[] = Array(queries.length).fill(0)
    queries.forEach(([threshold, count], qi) => {
      const y = minimize(
        preLen => {
          const sufLen = count - preLen
          return preSum[preLen] + 2 * threshold * sufLen - sufSum[sufLen]
        },
        0,
        count
      )[1]

      res[qi] = y
    })

    return res
  }
}
