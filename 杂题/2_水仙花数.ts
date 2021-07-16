// 水仙花数的定义：
// 一个 N 位非负整数，其各位数字的 N 次方和等于该数本身
const isNarcissisticNumber = (n: number) =>
  n ===
  n
    .toString()
    .split('')
    .map(str => parseInt(str, 10))
    .reduce((pre, cur, _, arr) => pre + Math.pow(cur, arr.length), 0)

const finaAllNarcissisticNumberByN = (n: number): number[] => {
  const res: number[] = []
  for (let index = Math.pow(10, n - 1); index < Math.pow(10, n); index++) {
    isNarcissisticNumber(index) && res.push(index)
  }
  return res
}

console.log(finaAllNarcissisticNumberByN(3))

export {}
