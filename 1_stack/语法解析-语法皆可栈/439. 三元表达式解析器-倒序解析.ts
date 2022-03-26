// 所包含的数字都只有一位数。

// 从后向前入栈，遇到问号就判断并出栈
function parseTernary(expression: string): string {
  const stack: string[] = []

  for (let i = expression.length - 1; ~i; i--) {
    if (expression[i] === '?') {
      const [num1, num2] = [stack.pop()!, stack.pop()!]
      const isTrue = expression[--i] === 'T'
      stack.push(isTrue ? num1 : num2)
    } else if (expression[i] !== ':') stack.push(expression[i])
  }

  return stack[0]
}

console.log(parseTernary('T?2:3'))
console.log(parseTernary('F?1:T?4:5'))
