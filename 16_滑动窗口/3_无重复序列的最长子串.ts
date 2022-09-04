/* eslint-disable @typescript-eslint/no-non-null-assertion */

const lengthOfLongestSubstring = (s: string): number => {
  const n = s.length
  const counter = new Map<string, number>()
  let left = 0
  let res = 0
  for (let right = 0; right < n; right++) {
    counter.set(s[right], (counter.get(s[right]) || 0) + 1)
    while (counter.get(s[right])! > 1) {
      counter.set(s[left], counter.get(s[left])! - 1)
      left++
    }
    res = Math.max(res, right - left + 1)
  }
  return res
}

console.log(lengthOfLongestSubstring('abcabcbb'))

export {}
