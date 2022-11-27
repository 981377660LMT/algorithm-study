/**
 * 1402. 做菜顺序
 * @param {number[]} satisfaction  1 <= n <= 500
 * @return {number}
 * 一个厨师收集了他 n 道菜的满意程度 satisfaction ，这个厨师做出每道菜的时间都是 1 单位时间。
 * 你可以按 任意 顺序安排做菜的顺序，你也可以选择放弃做某些菜来获得更大的总和。
 * 求time[i]*satisfaction[i] 最大值
 * @summary 正难则反
 * 我们可以考虑从满意度最高的菜开始选择，然后依次选择剩余的菜里满意度最高的，直到添加新菜不能使结果增加为止。
 */
function maxSatisfaction(satisfaction: number[]): number {
  satisfaction.sort((a, b) => b - a)

  let res = 0
  let preSum = 0
  for (const cur of satisfaction) {
    // 原来每个菜要增加一个单位的等待时间，因此原结果应增加s，而做这个菜也需一个单位时间，因此还应该增加num，也就是总共应该增加s+num
    // 每新加一个数，之前加过的所有数都会多加一遍
    const curSum = preSum + cur
    if (curSum > 0) {
      res += curSum
      preSum += cur
    } else {
      break
    }
  }

  return res
}

console.log(maxSatisfaction([-1, -8, 0, 5, -9]))

// 输入：satisfaction = [4,3,2]
// 输出：20
// 解释：按照原来顺序相反的时间做菜 (2*1 + 3*2 + 4*3 = 20)
// 输入：satisfaction = [-1,-8,0,5,-9]
// 输出：14
// 解释：去掉第二道和最后一道菜，最大的喜爱时间系数和为 (-1*1 + 0*2 + 5*3 = 14) 。每道菜都需要花费 1 单位时间完成。

export {}
