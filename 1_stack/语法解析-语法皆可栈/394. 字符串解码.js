/**
 * @param {string} s
 * @return {string}
 * @summary
 * 这道题和括号匹配一个思路
 * 遇到]就停下来，弹出栈  然后去除word部分和数字部分 重复再入栈
 */
const decodeString = function (s) {
  const isNumeric = str => !isNaN(parseFloat(str)) && isFinite(str)
  const stack = []

  for (const letter of s) {
    if (letter === ']') {
      let repeatStr = ''
      let repeatCount = ''

      while (stack.length && stack[stack.length - 1] !== '[') {
        repeatStr = stack.pop() + repeatStr
      }
      // pop 掉 "["
      stack.pop()

      while (stack.length && isNumeric(stack[stack.length - 1])) {
        repeatCount = stack.pop() + repeatCount
      }

      stack.push(repeatStr.repeat(Number(repeatCount)))
    } else {
      stack.push(letter)
    }
  }

  return stack.join('')
}

console.log(decodeString('3[a]2[bc]'))
// "aaabcbc"
console.log(decodeString('3[a2[c]]'))
// "accaccacc"
