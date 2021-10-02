function findContinuousSequence(target: number): number[][] {
  const res: number[][] = []
  let l = 1,
    r = 2,
    sum = 3

  while (l <= target >> 1) {
    if (sum === target) {
      res.push(range(l, r + 1))
      sum -= l
      l++
    } else if (sum < target) {
      r++
      sum += r
    } else {
      sum -= l
      l++
    }
  }

  return res

  function range(l: number, r: number) {
    const res: number[] = []
    for (let i = l; i < r; i++) {
      res.push(i)
    }
    return res
  }
}

console.log(findContinuousSequence(9))
// 输入一个正整数 target ，
// 输出所有和为 target 的连续正整数序列（至少含有两个数）。
// 输入：target = 9
// 输出：[[2,3,4],[4,5]]
