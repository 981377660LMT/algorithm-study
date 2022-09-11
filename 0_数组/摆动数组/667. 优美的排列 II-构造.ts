/* eslint-disable no-shadow */
// 构造题

/**
 * 请你构造一个答案列表 answer ，该列表应当包含从 1 到 n 的 n 个不同正整数,
 * !且相邻两个数差的绝对值有且仅有 k (k>=1)个不同整数(考虑摆动序列,震荡幅度越来越小)
 * 如果存在多种答案，只需返回其中 任意一种
 */
function constructArray(n: number, k: number): number[] {
  return [...solve1(1, k + 1), ...solve2(k + 2, n)]

  // !摆动数组，振幅越来越小
  function solve1(lower: number, upper: number): number[] {
    const res = [lower]
    const n = upper - lower
    let diff = upper - lower
    let flag = 1
    for (let _ = 0; _ < n; _++) {
      res.push(res[res.length - 1] + flag * diff)
      diff--
      flag *= -1
    }
    return res
  }

  // !相邻差值为1的数组
  function solve2(lower: number, upper: number): number[] {
    const res: number[] = []
    for (let i = lower; i <= upper; i++) {
      res.push(i)
    }
    return res
  }
}

console.log(constructArray(3, 1))
console.log(constructArray(3, 2))
// 1,2,3,4,5 假设4个不同
// 加4减3加2减1
