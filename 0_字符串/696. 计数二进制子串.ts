/**
 * @param {string} s
 * @return {number}
 * 计算具有相同数量 0 和 1 的非空（连续）子字符串的数量，
 * 并且这些子字符串中的所有 0 和所有 1 都是连续的。
 * @summary 考察交界处
 * 0011   [2,2]  =>2
 * 00011  [3,2]  =>2
 * 000111 [3,3]  =>3  交界处中心扩展
 */
const countBinarySubstrings = function (s: string): number {
  const matched = s.match(/([01])\1*/g)!.map(v => v.length)
  let res = 0
  matched.reduce((pre, cur) => {
    res += Math.min(pre, cur)
    return cur
  })
  return res
}

console.log(countBinarySubstrings('00110011'))
console.log(countBinarySubstrings('1010001'))

export {}
