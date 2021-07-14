/**
 * @description 特点
 * @description 应用场景
 * @description 时间复杂度 O(n)
 * @description 空间复杂度 O(n)
 */
const isValid = (str: string) => {
  const stack: string[] = []
  const leftBrackets = ['(', '{', '[']
  const rightBrackets = [')', '}', ']']

  for (const s of str) {
    if (leftBrackets.includes(s)) {
      stack.push(s)
    } else if (rightBrackets.includes(s)) {
      // 判断栈顶元素与当前右括号的关系，匹配则弹出左括号
      const last = stack.slice(-1)[0]
      if (leftBrackets.indexOf(last) === rightBrackets.indexOf(s)) {
        stack.pop()
      } else {
        // 不匹配则返回false
        return false
      }
    }
  }

  return stack.length === 0
}

console.log(isValid('()'))
console.log(isValid('()('))
console.log(isValid('[]{})'))
console.log(isValid('[{]}'))
console.log(isValid(''))

export {}
