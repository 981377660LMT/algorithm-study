// 双指针即可
const lengthOfLongestSubstring = (s: string): number => {
  let left = 0
  let right = 0
  let max = 0
  const map = new Map<string, boolean>()

  while (right <= s.length - 1) {
    if (!map.has(s[right])) {
      map.set(s[right], true)
      max = Math.max(map.size, max)
      right++
    } else {
      left++
      map.delete(s[left - 1])
    }
  }

  return max
}

console.log(lengthOfLongestSubstring('abcabcbb'))

export {}
