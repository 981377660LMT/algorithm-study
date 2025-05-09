const MAPPING: Record<string, (a: number, b: number) => number> = {
  '+': (a, b) => a + b,
  '-': (a, b) => a - b,
  '*': (a, b) => a * b,
  '/': (a, b) => ~~(a / b)
}
const OPT = new Set(['+', '-', '*', '/'])

// 逆波兰表达式(后缀表达式)求值
function evalRPN(tokens: string[]): number {
  const numStack: number[] = []
  for (const c of tokens) {
    if (!OPT.has(c)) {
      numStack.push(Number(c))
      continue
    }
    const [num2, num1] = [numStack.pop()!, numStack.pop()!]
    numStack.push(MAPPING[c](num1, num2))
  }
  return numStack[0]
}

console.log(evalRPN(['4', '13', '5', '/', '+']))

export {}
