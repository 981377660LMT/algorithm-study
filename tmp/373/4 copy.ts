import { SortedListFast } from '../../22_专题/离线查询/根号分治/SortedList/SortedListFast'

export {}

const INF = 2e9 // !超过int32使用2e15

// VOWELS = set(["a", "e", "i", "o", "u"])

// class Solution:
//     def beautifulSubstrings(self, s: str, k: int) -> int:
//         nums = [1 if c in VOWELS else 0 for c in s]
//         res = 0
//         for i in range(len(nums)):
//             v, c = 0, 0
//             for j in range(i, len(nums)):
//                 if nums[j] == 1:
//                     v += 1
//                 else:
//                     c += 1
//                 if v == c and v * c % k == 0:
//                     res += 1
//         return res

const VOWELS = new Set(['a', 'e', 'i', 'o', 'u'])
function beautifulSubstrings(s: string, k: number): number {
  const n = s.length
  const okLen = new Set<number>() // 允许的长度
  for (let c = 0; c < n; c++) {
    if ((c * c) % k === 0) {
      okLen.add(2 * c)
    }
  }
  const tmp = Array.from(okLen).sort((a, b) => a - b)
  const okLen2 = new Uint8Array(2 * n + 10)
  tmp.forEach(v => {
    okLen2[v] = 1
  })

  const nums = new Int8Array(n)
  for (let i = 0; i < n; i++) nums[i] = VOWELS.has(s[i]) ? 1 : -1

  let res = 0
  for (let i = 0; i < n; i++) {
    let diff = 0
    for (let j = i; j < n; j++) {
      diff += nums[j]
      if (diff === 0 && okLen2[j - i + 1]) {
        res++
      }
    }
  }
  return res
}
