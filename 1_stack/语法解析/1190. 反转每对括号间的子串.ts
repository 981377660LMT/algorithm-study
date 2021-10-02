/**
 *
 * @param s  0 <= s.length <= 2000
 * 按照从括号内到外的顺序，逐层反转每对匹配括号中的字符串，并返回最终的结果。
 */
function reverseParentheses(s: string): string {
  if (s.length <= 1) return s

  const indexStack: number[] = []
  const stringBuilder = s.split('')

  for (let i = 0; i < s.length; i++) {
    const char = s[i]
    if (char === '(') {
      indexStack.push(i)
    } else if (char === ')') {
      const start = indexStack.pop()!
      const end = i
      reverse(stringBuilder, start, end)
    }
  }

  return stringBuilder.filter(char => char !== '(' && char !== ')').join('')

  function reverse(arr: any[], i: number, j: number) {
    while (i < j) {
      ;[arr[i], arr[j]] = [arr[j], arr[i]]
      i++
      j--
    }
  }
}

console.log(reverseParentheses('(ed(et(oc))el)'))

export {}
