const str = '2(3)(1)'
const str2 = '4(2(3)(1))(6(5))(2(3)(1))(6(5))'
const match = str.match(/^-?\d+/g)
const bracket = str2.match(/(\(.*?\)){2}/g)
const bracket2 = str.match(/\(.*?\){2}/g)
console.log(match)
console.log(bracket)
console.log(bracket2)
export default 1

// 无法指定匹配括号的个数
// 需要用栈的思想
