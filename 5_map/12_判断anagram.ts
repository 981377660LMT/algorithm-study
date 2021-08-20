/**
 * @param {string} s
 * @param {string} t
 * @return {boolean}
 * 时间复杂度为O(n)，空间上因为定义是的一个常量大小的辅助数组，所以空间复杂度为O(1)。
 */
const isAnagram = function (s: string, t: string): boolean {
  if (s.length !== t.length) return false
  const resSet = Array<number>(26).fill(0)
  const base = 'a'.codePointAt(0)!

  for (const i of s) {
    resSet[i.codePointAt(0)! - base]++
  }

  for (const i of t) {
    if (resSet[i.codePointAt(0)! - base] === 0) return false
    resSet[i.codePointAt(0)! - base]--
  }

  return true
}
console.log(isAnagram('anagram', 'nagaram'))

export {}
