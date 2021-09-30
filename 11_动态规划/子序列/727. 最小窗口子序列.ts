import { bisectRight } from '../../9_排序和搜索/二分api/7_二分搜索寻找最插右入位置'

/**
 *
 * @param s1
 * @param s2
 * 找出 S1 中最短的（连续）子串 W ，使得 s2 是 W 的 子序列 。
 * @summary
 * 非滑动窗口解法
 */
function minWindow(s1: string, s2: string): string {
  if (s1.length < s2.length) return ''
  let res = ''

  const set = new Set(s2)
  // 记录s2中每个字符在s1中的出现位置 不断寻找下一个即可
  const charToIndex = new Map<string, number[]>()
  for (let i = 0; i < s1.length; i++) {
    if (set.has(s1[i])) {
      !charToIndex.has(s1[i]) && charToIndex.set(s1[i], [])
      charToIndex.get(s1[i])!.push(i)
    }
  }

  for (const first of charToIndex.get(s2[0]) || []) {
    let last = first
    let hit = 1 // 匹配的字符个数
    for (let i = 1; i < s2.length; i++) {
      const nextChar = s2[i]
      if (!charToIndex.has(nextChar)) continue
      const indexes = charToIndex.get(nextChar)!
      const nextIndex = bisectRight(indexes, last)
      if (nextIndex >= indexes.length) continue
      last = indexes[nextIndex]
      hit++
    }

    if (hit === s2.length && (res === '' || last - first + 1 < res.length)) {
      res = s1.slice(first, last + 1)
    }
  }

  return res
}

// 输入：
// S = "abcdebdde", T = "bde"
// 输出："bcde"
// 解释：
// "bcde" 是答案，因为它在相同长度的字符串 "bdde" 出现之前。
// "deb" 不是一个更短的答案，因为在窗口中必须按顺序出现 T 中的元素。
// console.log(minWindow('abcdebdde', 'bde'))
console.log(minWindow('cnhczmccqouqadqtmjjzl', 'mm'))
console.log(minWindow('jmeqksfrsdcmsiwvaovztaqenprpvnbstl', 'k'))
