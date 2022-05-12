// 如果 分隔符 是一个正则表达式，且包含捕获括号（capturing parentheses），
// 则每次匹配到 分隔符 时，捕获括号匹配的结果将会插入到返回的数组中
const s = '(1+(42+5+2)-3)+(6+8)'
const tokens = s.split(/([\(\)\+\-\*\/])/g).filter(Boolean)

// 包含分隔符
console.log(tokens)

// 不包含分隔符 python 也一样
console.log(s.split(/\+|\-|\*|\/|\(|\)/g))

export {}
