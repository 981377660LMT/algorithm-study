// 将 s 分割成 2 个字符串 p 和 q ，它们连接起来等于 s 且 p 和 q 中不同字符的数目相同。
// 1 <= s.length <= 10^5
function numSplits(s: string): number {
  let res = 0
  const allCounter = Array(26).fill(0)
  const visitedCounter = Array(26).fill(0)
  let allType = 0
  let visitedType = 0

  for (let i = 0; i < s.length; i++) {
    allCounter[s[i].codePointAt(0)! - 97]++
    if (allCounter[s[i].codePointAt(0)! - 97] === 1) allType++
  }

  for (let i = 0; i < s.length; i++) {
    const curCount = visitedCounter[s[i].codePointAt(0)! - 97]
    if (curCount === 1) visitedType++
    // 如果当前遍历到某个字母出现的次数是整个字符串该字母的总次数，剩下右部分不会再有该字母
    // 即右部分的不同字母数少一
    if (curCount === allCounter[s[i].codePointAt(0)! - 97]) allType--
    if (visitedType === allType) res++
  }

  return res
}

console.log(numSplits('aacaba'))
