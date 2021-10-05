function checkZeroOnes(s: string): boolean {
  return Math.max(...s.split('0').map(v => v.length)) > Math.max(...s.split('1').map(v => v.length))
}

console.log(checkZeroOnes('1101'))
// 如果字符串中由 1 组成的 最长 连续子字符串 严格长于 由 0 组成的 最长 连续子字符串，返回 true
// 由 1 组成的最长连续子字符串的长度是 2："1101"
// 由 0 组成的最长连续子字符串的长度是 1："1101"
// 由 1 组成的子字符串更长，故返回 true 。
