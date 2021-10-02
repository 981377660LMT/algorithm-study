// 考虑大数
function printNumbers(n: number): number[] {
  const res: string[] = []

  const bt = (index: number, digit: number, path: string[]) => {
    if (index === digit) return res.push(path.join(''))
    for (let next = 0; next <= 9; next++) {
      path.push(next.toString())
      bt(index + 1, digit, path)
      path.pop()
    }
  }

  for (let digit = 1; digit <= n; digit++) {
    for (let first = 1; first <= 9; first++) {
      bt(1, digit, [first.toString()])
    }
  }

  return res.map(Number)
}

console.log(printNumbers(2))
// 输入: n = 1
// 输出: [1,2,3,4,5,6,7,8,9]
