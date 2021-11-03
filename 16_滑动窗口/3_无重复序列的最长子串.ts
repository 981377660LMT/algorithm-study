// 双指针即可
const lengthOfLongestSubstring = (s: string): number => {
  let left = 0
  let right = 0
  let max = 0
  const visited = new Set<string>()

  while (right < s.length) {
    if (!visited.has(s[right])) {
      visited.add(s[right])
      right++
      max = Math.max(visited.size, max)
    } else {
      visited.delete(s[left])
      left++
    }
  }

  return max
}

console.log(lengthOfLongestSubstring('abcabcbb'))

export {}
