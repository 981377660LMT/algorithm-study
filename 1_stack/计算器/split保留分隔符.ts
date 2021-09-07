// 如果 separator 是一个正则表达式，且包含捕获括号（capturing parentheses），
// 则每次匹配到 separator 时，捕获括号匹配的结果将会插入到返回的数组中
const s = '(1+(4+5+2)-3)+(6+8)'
const tokens = s.split(/([\(\)\+\-\*\/])/g).filter(Boolean)
console.log(tokens)

export {}
