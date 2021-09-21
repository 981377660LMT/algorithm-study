/**
 * @param {string} s  s 由整数和算符 ('+', '-', '*', '/') 组成，中间由一些空格隔开
 * @return {number}
 */
var calculate = function (s: string): number {
  // 哨兵  类似"-1+3"或"1+(-1+2)"第一位是负数的表达式，会有问题 要加入0
  s = s.startsWith('-') ? '0' + s + '$' : s + '$'
  const tokens = s.split(/([\+\-\*\/\$])/g).filter(v => v.trim())

  // console.log(tokens)
  const opt = new Set(['+', '-', '*', '/', '$'])
  const numStack: number[] = []
  const optStack: string[] = []
  const evaluate = {
    '+': (a: number, b: number) => a + b,
    '-': (a: number, b: number) => a - b,
    '*': (a: number, b: number) => a * b,
    '/': (a: number, b: number) => ~~(a / b),
  } as Record<string, (a: number, b: number) => number>
  const getOperatorWeight = (operator: string) => {
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

  for (const token of tokens) {
    if (!opt.has(token)) {
      numStack.push(parseInt(token))
    } else {
      while (
        optStack.length &&
        getOperatorWeight(optStack[optStack.length - 1]) >= getOperatorWeight(token)
      ) {
        const [num2, num1] = [numStack.pop(), numStack.pop()] as [number, number]
        numStack.push(evaluate[optStack.pop()!](num1, num2))
      }
      optStack.push(token)
    }
  }

  return numStack[0]
}

console.log(calculate(' 3+5 / 2 '))
console.log(calculate('3+2*2'))

export {}
