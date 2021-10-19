/**
 * @param {number[]} temperatures
 * @return {number[]}
 * 请根据每日 气温 列表 temperatures ，请计算在每一天需要等几天才会有更高的温度。如果气温在这之后都不会升高，请在该位置用 0 来代替。
 */
const dailyTemperatures = function (temperatures: number[]): number[] {
  // temperatures.unshift(0)
  // temperatures.push(0)  可以不用
  const n = temperatures.length
  const res: number[] = Array(n).fill(0)
  const stack: number[] = []
  for (let i = 0; i < n; i++) {
    while (stack.length && temperatures[stack[stack.length - 1]] < temperatures[i]) {
      const tmp = stack.pop()!
      res[tmp] = i - tmp
    }

    stack.push(i)
  }

  return res
}

console.log(dailyTemperatures([73, 74, 75, 71, 69, 72, 76, 73]))
