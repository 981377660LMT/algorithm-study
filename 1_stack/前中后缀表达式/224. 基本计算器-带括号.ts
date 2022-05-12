const CAL: Record<string, (a: number, b: number) => number> = {
  '+': (a, b) => a + b,
  '-': (a, b) => a - b,
  '*': (a, b) => a * b,
  '/': (a, b) => ~~(a / b),
}

const WEIGHT: Record<string, number> = Object.fromEntries([
  ['(', NaN],
  ['+', 1],
  ['-', 1],
  ['*', 2],
  ['/', 2],
])

const OPT = new Set(['(', ')', '+', '-', '*', '/'])

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

  const numStack: number[] = []
  const optStack: string[] = []

  for (const char of tokens) {
    if (!OPT.has(char)) {
      numStack.push(Number(char))
    } else {
      if (char !== ')') {
        while (optStack.length && WEIGHT[optStack[optStack.length - 1]] >= WEIGHT[char]) {
          const [num2, num1] = [numStack.pop(), numStack.pop()] as [number, number]
          numStack.push(CAL[optStack.pop()!](num1, num2))
        }
        optStack.push(char)
      } else {
        while (optStack.length && optStack[optStack.length - 1] !== '(') {
          const [num2, num1] = [numStack.pop(), numStack.pop()] as [number, number]
          numStack.push(CAL[optStack.pop()!](num1, num2))
        }
        optStack.pop() // 弹出'('
      }
    }
  }

  return numStack[0]
}

console.log(calculate('(1+(4+5+2)-3)+(6+8)'))
// console.log(calculate('2-1+2'))
// console.log(calculate('2147483647'))
// console.log(calculate('-2+ 1'))
// console.log(calculate('1-(-2)'))
console.log(calculate('6-4/2'))

export default 1
