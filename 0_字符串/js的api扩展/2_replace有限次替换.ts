// 模拟python的replace方法
const replace = function (str: string, searchValue: string, replaceValue: string, limit = 1) {
  const searchValueRegexp = new RegExp(searchValue, 'g')
  return str.replace(searchValueRegexp, match => {
    return --limit >= 0 ? replaceValue : match
  })
}
// String.prototype.replace(searchValue, replaceValue)
console.log(replace('aaabbb', 'a', 'k'))
console.log(replace('aaabbb', 'a', 'k', 2))
console.log(replace('aaabbb', 'a', 'k', 3))
