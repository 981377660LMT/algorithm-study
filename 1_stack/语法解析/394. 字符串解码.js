/**
 * @param {string} s
 * @return {string}
 * @summary
 * 这道题和括号匹配一个思路
 * 遇到]就停下来，弹出栈  然后去除word部分和数字部分 重复再入栈
 */
const decodeString = function (s) {
  const isNumeric = str => !Number.isNaN(parseFloat(str)) && isFinite(str)
  const stack = []

  for (const letter of s) {
    if (letter !== ']') {
      stack.push(letter)
      continue
    }

    let word = ''
    let cur = stack.pop()
    while (cur !== '[') {
      word = cur + word
      cur = stack.pop()
    }

    let count = ''
    cur = stack.pop()
    while (isNumeric(cur)) {
      // 字符串表示数字保证正确
      count = cur + count
      cur = stack.pop()
    }

    stack.push(cur)
    stack.push(word.repeat(parseInt(count)))
  }

  return stack.join('')
}

console.log(decodeString('3[a]2[bc]'))
// "aaabcbc"
console.log(decodeString('3[a2[c]]'))
// "accaccacc"
