// 考虑大数
function printNumbers(n: number): number[] {
  const res: string[] = []

  for (let length = 1; length <= n; length++) {
    for (let first = 1; first <= 9; first++) {
      bt(1, length, [first.toString()])
    }
  }

  return res.map(Number)

  function bt(index: number, length: number, path: string[]): void {
    if (index === length) {
      res.push(path.join(''))
      return
    }

    for (let next = 0; next <= 9; next++) {
      path.push(next.toString())
      bt(index + 1, length, path)
      path.pop()
    }
  }
}

console.log(printNumbers(2))
// 输入: n = 1
// 输出: [1,2,3,4,5,6,7,8,9]
