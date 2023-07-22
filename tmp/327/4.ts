import { Trie } from '../../6_tree/前缀树trie/Trie'
import { BITArray } from '../../6_tree/树状数组/经典题/BIT'

export {}

const INF = 2e15
// 给你一个字符串 word 和一个字符串数组 forbidden 。

// 如果一个字符串不包含 forbidden 中的任何字符串，我们称这个字符串是 合法 的。

// 请你返回字符串 word 的一个 最长合法子字符串 的长度。

// 子字符串 指的是一个字符串中一段连续的字符，它可以为空。
// 如果10可以，那么10以上一定可以

function longestValidSubstring(word: string, forbidden: string[]): number {
  const n = word.length
  const allBad = new Set<string>(forbidden)
  const bads = [...allBad].sort((a, b) => a.length - b.length).slice(0, 1e3)
  let res = 0
  let left = 0
  for (let right = 0; right < n; right++) {
    while (left <= right && !check(left, right)) {
      left++
    }
    res = Math.max(res, right - left + 1)
  }
  return res

  function check(left: number, right: number): boolean {
    if (left > right) return true
    const cur = word.slice(left, Math.min(right + 1, 1e4))
    for (let i = 0; i < bads.length; i++) {
      const b = bads[i]
      const index = cur.indexOf(b)
      if (index !== -1 && index + b.length - 1 <= right) return false
    }
    return true
  }
}

if (require.main === module) {
  const n = 1e5
  const m = 1e5
  const word = Array.from({ length: n }, () =>
    String.fromCharCode(Math.floor(Math.random() * 26) + 97)
  ).join('')
  const forbidden = Array.from({ length: m }, () => {
    const len = Math.floor(Math.random() * 10) + 1
    return Array.from({ length: len }, () =>
      String.fromCharCode(Math.floor(Math.random() * 26) + 97)
    ).join('')
  })

  console.time('as')
  console.log(longestValidSubstring(word, forbidden))
  console.timeEnd('as')
}
