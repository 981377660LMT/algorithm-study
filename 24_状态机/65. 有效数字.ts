type Token = 'SIGN' | 'DIGIT' | 'DOT' | 'EXP' | 'NONE'

/**
 * @param {string} s
 * @return {boolean}
 * 给你一个字符串 s ，如果 s 是一个 有效数字 ，请返回 true 。
 * @summary
 * 状态转移规则很繁琐
 */
const isNumber = function (s: string): boolean {
  const states = {
    start: { SIGN: 'sign1', DIGIT: 'digit1', DOT: 'dot1' },
    sign1: { DIGIT: 'digit1', DOT: 'dot1' }, // 非exp后的符号
    sign2: { DIGIT: 'D' }, // exp后的符号
    digit1: { DIGIT: 'digit1', DOT: 'dot2', EXP: 'exp', END: '' }, // 前面带符号
    digit2: { DIGIT: 'digit2', EXP: 'exp', END: '' },
    dot1: { DIGIT: 'digit2' }, // 前面没数字
    dot2: { DIGIT: 'digit2', EXP: 'exp', END: '' }, // 前面有数字
    exp: { SIGN: 'sign2', DIGIT: 'D' },
    D: { DIGIT: 'D', END: '' },
  } as Record<string, Partial<Record<Token, string>>>

  const isDigit = (x: any) => !isNaN(parseFloat(x)) && isFinite(x)

  const getTokenType = (str: string): Token => {
    if (str === '.') return 'DOT'
    else if (['+', '-'].includes(str)) return 'SIGN'
    else if (['E', 'e'].includes(str)) return 'EXP'
    else if (isDigit(str)) return 'DIGIT'
    return 'NONE'
  }

  let state = 'start'
  for (const char of s) {
    const tokenType = getTokenType(char)
    if (!(tokenType in states[state])) return false
    // 状态转移
    state = states[state][tokenType]!
  }

  // 最后是结束的状态
  return !!state && 'END' in states[state]
}

console.log(isNumber('.1'))
console.log(isNumber('e'))
console.log(isNumber('1 '))
