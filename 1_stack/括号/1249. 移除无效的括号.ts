/**
 * @param {string} s
 * @return {string}
 * 你需要从字符串中删除最少数目的 '(' 或者 ')' （可以删除任意位置的括号)，使得剩下的「括号字符串」有效。
 * 请返回任意一个合法字符串。
 * @summary
 * 统计哪些位置的括号不合法
 */
function minRemoveToMakeValid(s: string): string {
  const chars = s.split('')
  const invalid: number[] = []

  for (const [index, char] of chars.entries()) {
    if (char === '(') {
      invalid.push(index)
    } else if (char === ')') {
      if (invalid.length > 0) invalid.pop()
      else chars[index] = '' // 删不合法的右括号
    }
  }

  for (const index of invalid) {
    chars[index] = ''
  }

  return chars.join('')
}

console.log(minRemoveToMakeValid('lee(t(c)o)de)'))
