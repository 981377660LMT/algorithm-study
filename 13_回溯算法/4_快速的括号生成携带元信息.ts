function generateParenthesis(n: number): string[] {
  const res: string[] = []
  bt([], n, n)
  return res

  // 在bt中传递元信息参数加速验证
  function bt(stringBuilder: string[], leftRemaining: number, rightRemaining: number): void {
    if (rightRemaining < leftRemaining) return
    if (leftRemaining === 0 && rightRemaining === 0) {
      res.push(stringBuilder.join(''))
      return
    }

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
}

console.log(generateParenthesis(3))

export {}
