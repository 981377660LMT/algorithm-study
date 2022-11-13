// class Solution:
//     def maxPalindromes(self, s: str, k: int) -> int:
//         def helper(left, right):
//             """中心扩展法求所有回文子串"""
//             while left >= 0 and right < len(s) and s[left] == s[right]:
//                 left -= 1
//                 right += 1
//                 if right - left - 1 >= k:
//                     intervals.append((left + 1, right - 1))

//         n = len(s)
//         intervals = []
//         for i in range(n):
//             helper(i, i)
//             helper(i, i + 1)

//         intervals.sort(key=lambda x: x[1])
//         res = 0
//         preEnd = -1
//         for start, end in intervals:
//             if start > preEnd:
//                 res += 1
//                 preEnd = end
//         return res

// !注意会2000*2000会爆空间 大空间不要使用javascript
function maxPalindromes(s: string, k: number): number {
  const n = s.length
  const intervals: [number, number][] = []
  const expand = (left: number, right: number) => {
    while (left >= 0 && right < n && s[left] === s[right]) {
      left--
      right++
      if (right - left - 1 >= k) {
        intervals.push([left + 1, right - 1])
      }
    }
  }

  for (let i = 0; i < n; i++) {
    expand(i, i)
    expand(i, i + 1)
  }

  intervals.sort((a, b) => a[1] - b[1])
  let res = 0
  let preEnd = -1
  for (const [start, end] of intervals) {
    if (start > preEnd) {
      res++
      preEnd = end
    }
  }

  return res
}

export {}
