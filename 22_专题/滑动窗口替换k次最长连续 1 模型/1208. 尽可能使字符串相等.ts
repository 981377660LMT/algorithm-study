/**
 * @param {string} s
 * @param {string} t
 * @param {number} maxCost
 * @return {number}
 * 424
 * @description
 * 将 s 中的第 i 个字符变到 t 中的第 i 个字符需要 |s[i] - t[i]| 的开销（开销可能为 0），也就是两个字符的 ASCII 码值的差的绝对值。
 * 如果你可以将 s 的子字符串转化为它在 t 中对应的子字符串，则返回可以转化的最大长度。
 */
const equalSubstring = function (s: string, t: string, maxCost: number): number {
  let res = 0
  let left = 0
  let cost = 0

  for (let right = 0; right < s.length; right++) {
    cost += Math.abs(s[right].codePointAt(0)! - t[right].codePointAt(0)!)

    while (cost > maxCost) {
      cost -= Math.abs(s[left].codePointAt(0)! - t[left].codePointAt(0)!)
      left++
    }

    res = Math.max(res, right - left + 1)
  }

  return res
}

console.log(equalSubstring('abcd', 'bcdf', 3))
