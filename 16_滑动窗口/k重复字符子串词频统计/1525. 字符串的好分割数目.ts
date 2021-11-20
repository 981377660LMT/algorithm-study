// 将 s 分割成 2 个字符串 p 和 q ，它们连接起来等于 s 且 p 和 q 中不同字符的数目相同。
// 1 <= s.length <= 10^5
function numSplits(s: string): number {
  let res = 0
  const rightCounter = Array(26).fill(0)
  const leftCounter = Array(26).fill(0)
  let rightType = 0
  let leftType = 0

  for (let i = 0; i < s.length; i++) {
    rightCounter[s[i].codePointAt(0)! - 97]++
    if (rightCounter[s[i].codePointAt(0)! - 97] === 1) rightType++
  }

  for (let i = 0; i < s.length; i++) {
    leftCounter[s[i].codePointAt(0)! - 97]++
    if (leftCounter[s[i].codePointAt(0)! - 97] === 1) leftType++
    // 如果当前遍历到某个字母出现的次数是整个字符串该字母的总次数，剩下右部分不会再有该字母
    // 即右部分的不同字母数少一
    if (leftCounter[s[i].codePointAt(0)! - 97] === rightCounter[s[i].codePointAt(0)! - 97])
      rightType--
    if (leftType === rightType) res++
  }

  return res
}

console.log(numSplits('aacaba'))

export {}
