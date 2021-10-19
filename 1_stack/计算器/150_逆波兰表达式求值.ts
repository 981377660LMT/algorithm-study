/**
 * @param {string[]} tokens
 * @return {number}
 */
const evalRPN = function (tokens: string[]): number {
  const operators = new Set(['+', '-', '*', '/'])
  const stack: (string | number)[] = []
  const evaluate = {
    '+': (a: number, b: number) => a + b,
    '-': (a: number, b: number) => a - b,
    '*': (a: number, b: number) => a * b,
    '/': (a: number, b: number) => ~~(a / b),
  } as Record<string, (a: number, b: number) => number>

  for (const char of tokens) {
    if (!operators.has(char)) {
      stack.push(Number(char))
    } else {
      const [num2, num1] = [stack.pop(), stack.pop()] as [number, number]
      stack.push(evaluate[char](num1, num2))
    }
  }

  return stack[0] as number
}

console.log(evalRPN(['4', '13', '5', '/', '+']))
