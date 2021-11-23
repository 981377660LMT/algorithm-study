type Binary = (a: number, b: number) => number
/**
 * @param {string} s  s 由数字、'+'、'-'、'('、')'、和 ' ' 组成
 * @return {number}
 */
const calculate = function (s: string): number {
  // 哨兵  类似"-1+3"或"1+(-1+2)"第一位是负数的表达式，会有问题 要加入0
  // 加入')'保证所有操作符最后出栈
  s = s.startsWith('-') ? '0' + s + ')' : s + ')' // "-1 - (-1)" => "0-1 - (-1)"
  s = s.replace(/\(\-/g, '(0-') // "1 - (-1)" => "1 - (0-1)"
  s = s.replace(/\(\+/g, '(0+') // "1 - (+1)" => "1 - (0+1)"
  const tokens = s.split(/([\(\)\+\-\*\/])/g).filter(v => v.trim())
  const opt = new Set(['+', '-', '(', ')', '*', '/'])
  const numStack: number[] = []
  const optStack: string[] = []
  const evaluate: Record<string, Binary> = {
    '+': (a, b) => a + b,
    '-': (a, b) => a - b,
    '*': (a, b) => a * b,
    '/': (a, b) => ~~(a / b),
  }
  const getOperatorWeight = (operator: string) => {
    switch (operator) {
      case '+':
      case '-':
        return 0
      case '*':
      case '/':
        return 1
      case '(':
      case ')':
        return NaN
      default:
        throw new Error('Invalid Operator')
    }
  }

  for (const token of tokens) {
    if (!opt.has(token)) {
      numStack.push(parseInt(token))
    } else {
      if (token !== ')') {
        while (
          optStack.length &&
          getOperatorWeight(optStack[optStack.length - 1]) >= getOperatorWeight(token)
        ) {
          const [num2, num1] = [numStack.pop(), numStack.pop()] as [number, number]
          numStack.push(evaluate[optStack.pop()!](num1, num2))
        }
        optStack.push(token)
      } else {
        while (optStack.length && optStack[optStack.length - 1] !== '(') {
          const [num2, num1] = [numStack.pop(), numStack.pop()] as [number, number]
          numStack.push(evaluate[optStack.pop()!](num1, num2))
        }
        optStack.pop() // 弹出'('
      }
    }
  }

  return numStack[0]
}

// console.log(calculate('(1+(4+5+2)-3)+(6+8)'))
// console.log(calculate('2-1+2'))
// console.log(calculate('2147483647'))
// console.log(calculate('-2+ 1'))
// console.log(calculate('1-(-2)'))
console.log(calculate('6-4/2'))

export default 1
