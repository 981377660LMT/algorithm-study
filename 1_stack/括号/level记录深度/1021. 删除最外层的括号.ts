// 1021. 删除最外层的括号

function removeOuterParentheses(s: string): string {
  const sb: string[] = []
  let level = 0

  for (const bracket of s) {
    if (bracket === '(') {
      level++
      if (level > 1) sb.push(bracket)
    } else {
      level--
      if (level > 0) sb.push(bracket)
    }
  }

  return sb.join('')
}

console.log(removeOuterParentheses('(()())(())'))
// 输出："()()()"
console.log(removeOuterParentheses('()()'))
// 输出：""
