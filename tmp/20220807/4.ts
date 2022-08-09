function longestIdealString(s: string, k: number): number {
  const n = s.length
  const memo = new Int32Array(30 * n).fill(-1)

  const dfs = (index: number, pre: number): number => {
    if (index === n) return 0
    const hash = index * 29 + pre
    if (memo[hash] !== -1) return memo[hash]
    if (pre === 0) {
      const cand1 = dfs(index + 1, pre)
      const cand2 = dfs(index + 1, s.charCodeAt(index) - 96) + 1
      const res = Math.max(cand1, cand2)
      memo[hash] = res
      return res
    }

    let res = dfs(index + 1, pre)
    if (Math.abs(pre - (s.charCodeAt(index) - 96)) <= k) {
      const cand2 = dfs(index + 1, s.charCodeAt(index) - 96) + 1
      res = Math.max(res, cand2)
    }
    memo[hash] = res
    return res
  }

  return dfs(0, 0)
}

export {}

if (require.main === module) {
  // s = "abcd", k = 3
  console.log(longestIdealString('abcd', 3))
  console.log(longestIdealString('acfgbd', 2))
}
