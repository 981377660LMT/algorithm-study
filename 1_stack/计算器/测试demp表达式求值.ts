type ArrayElement<ArrayType extends readonly unknown[]> =
  ArrayType extends readonly (infer ElementType)[] ? ElementType : never

const evalTokens = (tokens: any[]) => {
  // 哨兵，清空操作符栈
  tokens.push('$')
  const operators = ['+', '-', '*', '/', '$'] as const
  const set = new Set(['+', '-', '*', '/', '$'])
  const numStack: number[] = []
  const operatorStack: ArrayElement<typeof operators>[] = []
  const evaluate = {
    '+': (a: number, b: number) => a + b,
    '-': (a: number, b: number) => a - b,
    '*': (a: number, b: number) => a * b,
    '/': (a: number, b: number) => ~~(a / b),
  } as Record<string, (a: number, b: number) => number>
  const getOperatorWeight = (operator: ArrayElement<typeof operators>) => {
    switch (operator) {
      case '+':
      case '-':
        return 0
      case '*':
      case '/':
        return 1
      case '$':
        return -1
      default:
        throw new Error('Not Valid Operator')
    }
  }

  for (const letter of tokens) {
    if (!set.has(letter)) {
      numStack.push(parseInt(letter))
    } else {
      while (
        operatorStack.length &&
        getOperatorWeight(operatorStack[operatorStack.length - 1]) >= getOperatorWeight(letter)
      ) {
        const [num2, num1] = [numStack.pop(), numStack.pop()] as [number, number]
        numStack.push(evaluate[operatorStack.pop()!](num1, num2))
      }
      operatorStack.push(letter)
    }
  }

  return numStack[0]
}

console.log(evalTokens(['34', '+', '13', '*', '9', '+', '44', '-', '12', '/', '3']))
export default 1
