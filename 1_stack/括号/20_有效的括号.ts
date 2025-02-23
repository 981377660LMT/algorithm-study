function isValid(str: string) {
  const stack: string[] = []
  const mp = new Map<string, string>()
  mp.set('(', ')').set('{', '}').set('[', ']') // 左括号 => 右括号

  for (const s of str) {
    if (mp.has(s)) {
      stack.push(s)
    } else {
      // 判断栈顶元素与当前右括号的关系，匹配则弹出左括号
      if (stack.length === 0) return false
      const last = stack[stack.length - 1]
      if (mp.get(last) === s) {
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
