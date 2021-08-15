/**
 * @param {number[]} nums
 * @param {number} k
 * @return {number}
 * @description
 * 有一个仅包含’a’和’b’两种字符的字符串 s，长度为 n，每次操作可以把一个字符做一次转换（把一个’a’设置为’b’，或者把一个’b’置成’a’)；
 * 但是操作的次数有上限 k，问在有限的操作数范围内，能够得到最大连续的相同字符的子串的长度是多少。
 * 实际上题目中是求连续 a 或者 b 的长度。看到连续，大家也应该有滑动窗口的敏感度， 别管行不行， 想到总该有的。
 */
const longestOnes = function (s: string, k: number): number {
  const helper = (onlyValue: string, s: string, k: number): number => {
    let l = 0
    let r = 0
    let res = 0
    while (r < s.length) {
      // console.log(l, r)
      if (s[r] !== onlyValue) k--
      while (k < 0 && l < r) {
        l++
        if (s[l - 1] === onlyValue) k++
      }
      console.log(res, l, r)
      res = Math.max(res, r - l + 1)
      r++
    }
    return res
  }
  return Math.max(helper('1', s, k), helper('0', s, k))
}

console.log(longestOnes([0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1, 1].join(''), 3))

export default 1
