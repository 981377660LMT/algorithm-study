import { fix } from './最长连续 1 模型'

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
  return Math.max(fix(s, '1', k), fix(s, '0', k))
}

console.log(longestOnes([0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1, 1].join(''), 3))

export default 1
