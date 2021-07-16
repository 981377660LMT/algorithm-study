const generateParenthesis = (n: number): string[] => {
  const res: string[] = []

  // 在bt中传递元信息参数加速验证
  const bt = (path: string, leftRemaining: number, rightRemaining: number) => {
    if (leftRemaining === 0 && rightRemaining === 0) return res.push(path)

    if (rightRemaining < leftRemaining) return

    leftRemaining > 0 && bt(path + '(', leftRemaining - 1, rightRemaining)
    rightRemaining > 0 && bt(path + ')', leftRemaining, rightRemaining - 1)
  }
  bt('', n, n)

  return res
}

console.log(generateParenthesis(3))

export {}
