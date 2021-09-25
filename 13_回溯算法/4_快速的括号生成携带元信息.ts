const generateParenthesis = (n: number): string[] => {
  const res: string[] = []

  // 在bt中传递元信息参数加速验证
  const bt = (stringBuilder: string[], leftRemaining: number, rightRemaining: number) => {
    if (leftRemaining === 0 && rightRemaining === 0) return res.push(stringBuilder.join(''))

    if (rightRemaining < leftRemaining) return

    if (leftRemaining > 0) {
      stringBuilder.push('(')
      bt(stringBuilder, leftRemaining - 1, rightRemaining)
      stringBuilder.pop()
    }
    if (rightRemaining > 0) {
      stringBuilder.push(')')
      bt(stringBuilder, leftRemaining, rightRemaining - 1)
      stringBuilder.pop()
    }
  }
  bt([], n, n)

  return res
}

console.log(generateParenthesis(3))

export {}
