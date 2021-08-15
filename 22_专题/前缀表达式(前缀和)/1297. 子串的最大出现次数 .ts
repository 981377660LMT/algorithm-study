/**
 * @param {string} s
 * @param {number} maxLetters
 * @param {number} minSize
 * @param {number} maxSize
 * @return {number}
 * @description 请你返回满足以下条件且出现次数最大的 任意 子串的出现次数
 * 子串中不同字母的数目必须小于等于 maxLetters 。
   子串的长度必须大于等于 minSize 且小于等于 maxSize
   @summary 这道题的 maxSize 是无用信息， 并不需要用到(因为子串越短越好)。
   固定大小的滑动窗口
 */
const maxFreq = function (s: string, maxLetters: number, minSize: number, maxSize: number): number {
  const mapper = new Map<string, number>()

  const isValid = (str: string) => {
    const set = new Set<string>()
    for (const letter of str) set.add(letter)
    return set.size <= maxLetters
  }

  for (let i = 0; i <= s.length - minSize; i++) {
    const subStr = s.slice(i, i + minSize)
    if (!isValid(subStr)) continue
    mapper.set(subStr, (mapper.get(subStr) || 0) + 1)
  }
  console.log(Math.max())
  return Math.max(...mapper.values(), 0)
}

console.log(maxFreq('aababcaab', 2, 3, 4))
console.log(maxFreq('abcde', 2, 3, 3))
// 输出：2
// 解释：子串 "aab" 在原字符串中出现了 2 次。
// 它满足所有的要求：2 个不同的字母，长度为 3 （在 minSize 和 maxSize 范围内）。
