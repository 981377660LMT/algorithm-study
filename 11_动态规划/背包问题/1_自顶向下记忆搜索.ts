// 类比带cache的斐波那契数列计算
const recursionKnapsack = (
  stuffVolume: number[],
  stuffValue: number[],
  knapsackVolume: number
): number => {
  const num = stuffVolume.length
  // 记忆化搜索存储，保存"steps#volume"-value的值
  const memo = new Map<string, number>()
  const recur = (
    curVolumn: number,
    steps: number,
    sVol = stuffVolume,
    sVal = stuffValue,
    sum = 0
  ): number => {
    if (steps >= num || curVolumn <= 0) return sum
    if (memo.has(`${steps}#${curVolumn}`)) return memo.get(`${steps}#${curVolumn}`)!

    let putThisOrNot = recur(curVolumn, steps + 1, sVol, sVal, sum)
    if (curVolumn >= sVol[steps]) {
      const putThis = recur(curVolumn - sVol[steps], steps + 1, sVol, sVal, sum + sVal[steps])
      putThisOrNot = Math.max(putThisOrNot, putThis)
    }

    memo.set(`${steps}#${curVolumn}`, putThisOrNot)
    return putThisOrNot
  }

  return recur(knapsackVolume, 0)
}

console.dir(recursionKnapsack([1, 2, 3], [6, 10, 12], 5), { depth: null })

export {}
