// console.log(decodeString('3[a]2[bc]'))
//  "aaabcbc"
// console.log(decodeString('3[a2[c]]'))
//  "accaccacc"

function isNumeric(str: string) {
  const num = parseFloat(str)
  return !Number.isNaN(num) && Number.isFinite(num)
}

/**
 * 这道题和括号匹配一个思路
 * 遇到]就停下来，弹出栈  然后去除word部分和数字部分 重复再入栈
 */
function decodeString(s: string) {
  const stack: string[] = []

  for (const char of s) {
    if (char !== ']') {
      stack.push(char)
    } else {
      let word = ''
      let count = ''

      // !1. 找字符串
      while (stack.length && stack[stack.length - 1] !== '[') word = stack.pop() + word
      stack.pop() // pop 掉 "["

      // !2.找重复次数
      while (stack.length && isNumeric(stack[stack.length - 1])) count = stack.pop() + count
      stack.push(word.repeat(Number(count)))
    }
  }

  return stack.join('')
}

console.log(decodeString('3[a]2[bc]'))
console.log(decodeString('3[ab]'))
// "aaabcbc"
console.log(decodeString('3[a2[c]]'))
// "accaccacc"

export {}
