import { SortedListFast } from '../../22_专题/离线查询/根号分治/SortedList/SortedListFast'

export {}

const INF = 2e9 // !超过int32使用2e15

// class Solution:
//     def lexicographicallySmallestArray(self, nums: List[int], limit: int) -> List[int]:
//         # 维护一个区间最小值
//         res = []
//         sl = SortedList(nums)
//         for i, v in enumerate(nums):
//             # >= v - limit 的后继
//             pos = sl.bisect_left(v - limit)
//             cur = sl[pos]
//             res.append(cur)
//             sl.remove(cur)

//         return res
// VOWELS = set(["a", "e", "i", "o", "u"])

// class Solution:
//     def beautifulSubstrings(self, s: str, k: int) -> int:
//         n = len(s)
//         okLen = set()  # 允许的长度
//         for c in range(n):
//             if c * c % k == 0:
//                 okLen.add(2 * c)
//         tmp = sorted(okLen)
//         okLen = [False] * (2 * n + 10)
//         for v in tmp:
//             okLen[v] = True

//         nums = [1 if c in VOWELS else -1 for c in s]
//         res = 0
//         preSum = defaultdict(list)
//         preSum[0].append(-1)
//         curSum = 0
//         for i, v in enumerate(nums):
//             curSum += v
//             for j in preSum[curSum]:
//                 if okLen[i - j]:
//                     res += 1
//             preSum[curSum].append(i)
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
  for (let i = 0; i < n; i++) {
    if (VOWELS.has(s[i])) {
      nums[i] = 1
    } else {
      nums[i] = -1
    }
  }

  let res = 0
  const preSum = new Map<number, number[]>()
  preSum.set(0, [-1])
  let curSum = 0
  for (let i = 0; i < n; i++) {
    curSum += nums[i]
    const arr = preSum.get(curSum) || []
    for (let j = 0; j < arr.length; j++) {
      res += okLen2[i - arr[j]]
    }
    arr.push(i)
    preSum.set(curSum, arr)
  }
  return res
}
