// 给你一个平衡字符串 s，请你将它分割成尽可能多的平衡字符串。
// 返回可以通过分割得到的平衡字符串的 最大数量 。
function balancedStringSplit(s: string): number {
  let res = 0
  let sum = 0

  for (const char of s) {
    if (char === 'R') sum++
    else sum--
    if (sum === 0) res++
  }

  return res
}

console.log(balancedStringSplit('RLLLLRRRLR'))
// 输出：3
// 解释：s 可以分割为 "RL"、"LLLRRR"、"LR" ，每个子字符串中都包含相同数量的 'L' 和 'R' 。
