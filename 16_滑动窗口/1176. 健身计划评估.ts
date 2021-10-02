function dietPlanPerformance(calories: number[], k: number, lower: number, upper: number): number {
  let res = 0
  let sum = calories.slice(0, k).reduce((pre, cur) => pre + cur, 0)
  if (sum > upper) res++
  else if (sum < lower) res--

  for (let r = k; r < calories.length; r++) {
    sum += calories[r] - calories[r - k]
    if (sum > upper) res++
    else if (sum < lower) res--
  }

  return res
}
// 「这一天以及之后的连续几天」 （共 k 天）内消耗的总卡路里 T：

// 如果 T < lower，那么这份计划相对糟糕，并失去 1 分；
// 如果 T > upper，那么这份计划相对优秀，并获得 1 分；
// 否则，这份计划普普通通，分值不做变动。
