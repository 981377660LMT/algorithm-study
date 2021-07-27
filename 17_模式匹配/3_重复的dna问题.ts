/**
 * @param {string} s
 * @return {string[]}
 * @description 编写一个函数来找出所有目标子串，目标子串的长度为 10，且在 DNA 字符串 s 中出现次数超过一次
 * 固定长度的序列使用滑动窗口
 */
const findRepeatedDnaSequences = function (s: string): string[] {
  if (s.length <= 10) return []
  let cur = s.slice(0, 10)
  const visited = new Set<string>([cur])
  const res = new Set<string>()
  let i = 0
  let j = 9

  while (j < s.length) {
    i++
    j++
    cur = cur.slice(1) + s[j]
    if (visited.has(cur)) {
      res.add(cur)
    } else {
      visited.add(cur)
    }
  }

  return [...res]
}

console.log(findRepeatedDnaSequences('AAAAACCCCCAAAAACCCCCCAAAAAGGGTTT'))

export {}
