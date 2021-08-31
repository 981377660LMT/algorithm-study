// 双指针即可
const lengthOfLongestSubstring = (s: string): number => {
  let left = 0
  let right = 0
  let max = 0
  const set = new Set<string>()

  while (right <= s.length - 1) {
    if (!set.has(s[right])) {
      set.add(s[right])
      max = Math.max(set.size, max)
      right++
    } else {
      left++
      set.delete(s[left - 1])
    }
  }

  return max
}

console.log(lengthOfLongestSubstring('abcabcbb'))

export {}
