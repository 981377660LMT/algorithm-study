// class Solution:
//     def deleteString(self, s: str) -> int:
//         @lru_cache(None)
//         def dfs(index: int) -> int:
//             if index == S:
//                 return 0
//             if index == S - 1:
//                 return 1

//             remain = S - index
//             res = 1
//             for i in range(1, remain // 2 + 1):
//                 if s[index : i + index] == s[i + index : i + i + index]:
//                     cand = dfs(i + index) + 1
//                     res = cand if cand > res else res
//             return res

//         S = len(s)
//         res = dfs(0)
//         dfs.cache_clear()
//         return res

function deleteString(s: string): number {
  const S = s.length
  const memo = new Int16Array(S + 10).fill(-1)
  const res = dfs(0)
  return res
  function dfs(index: number): number {
    if (index === S) return 0
    if (index === S - 1) return 1
    if (memo[index] !== -1) return memo[index]
    let res = 1
    for (let i = 1; i <= (S - index) / 2; i++) {
      if (s.slice(index, i + index) === s.slice(i + index, i + i + index)) {
        const cand = dfs(i + index) + 1
        res = cand > res ? cand : res
      }
    }
    memo[index] = res
    return res
  }
}

export {}
