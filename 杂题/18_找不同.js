/**
 * @param {string} s
 * @param {string} t
 * @return {character}
 */
var findTheDifference = function (s, t) {
  const sum1 = s.split('').reduce((acc, cur) => acc + cur.charCodeAt(0), 0)
  const sum2 = t.split('').reduce((acc, cur) => acc + cur.charCodeAt(0), 0)
  return String.fromCodePoint(sum2 - sum1)
}

console.log(findTheDifference('as', 'abs'))
console.log(String.fromCodePoint(97, 98, 99))
console.log(String.fromCharCode(97, 98, 99))
// String.fromCharCode()最大支持16位的数字，而且ES中很早就开始支持
// ，兼容性好。而String.fromCodePoint()可以多达21位数字，
// 是ES6新加入的方法，是为了能够处理所有合法的
