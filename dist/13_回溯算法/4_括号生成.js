'use strict'
Object.defineProperty(exports, '__esModule', { value: true })
const generateParenthesis = n => {
  const res = []
  const parenthesis = ['(', ')']
  const isValidPath = path => {
    const stack = []
    for (const letter of path) {
      if (letter === '(') {
        stack.push(letter)
      } else if (letter === ')') {
        const head = stack.pop()
        if (head !== '(') return false
      }
    }
    return stack.length === 0
  }
  const bt = path => {
    if (path.length === 6) {
      if (isValidPath(path)) res.push(path)
      return
    }
    for (const p of parenthesis) {
      path += p
      bt(path)
    }
  }
  bt('')
  return res
}
console.log(generateParenthesis(3))
