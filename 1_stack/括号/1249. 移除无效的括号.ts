/**
 * @param {string} s
 * @return {string}
 * 你需要从字符串中删除最少数目的 '(' 或者 ')' （可以删除任意位置的括号)，使得剩下的「括号字符串」有效。
 * 请返回任意一个合法字符串。
 * @summary
 * 统计哪些位置的括号不合法
 */
var minRemoveToMakeValid = function (s: string): string {
  const chars = s.split('')
  const invalidIndexStack: number[] = []

  for (let i = 0; i < chars.length; i++) {
    const char = chars[i]
    if (char === '(') {
      invalidIndexStack.push(i)
    } else if (char === ')') {
      if (invalidIndexStack.length) invalidIndexStack.pop()
      else chars[i] = '' // 删不合法的右括号
    }
  }

  for (const index of invalidIndexStack) {
    chars[index] = ''
  }

  return chars.join('')
}

console.log(minRemoveToMakeValid('lee(t(c)o)de)'))
