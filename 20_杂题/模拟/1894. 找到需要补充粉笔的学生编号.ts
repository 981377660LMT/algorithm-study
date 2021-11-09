// 请你返回需要 补充 粉笔的学生 编号 。
function chalkReplacer(chalk: number[], k: number): number {
  const total = chalk.reduce((pre, cur) => pre + cur, 0)
  const remain = k % total
  let sum = 0

  for (const [index, value] of chalk.entries()) {
    sum += value
    if (sum > remain) return index
  }

  return -1
}

console.log(chalkReplacer([5, 1, 5], 22))
