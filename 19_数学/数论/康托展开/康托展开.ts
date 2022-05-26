// 可以记忆化dfs处理
const factorial = [1, 1, 2, 6, 24, 120, 720, 5040, 40320, 362880]

// 康托展开
function calRankByPerm(perm: number | string, n: number): number {
  const str = perm.toString()
  const counter = new Map<string, number>()
  for (const char of str) counter.set(char, (counter.get(char) ?? 0) + 1)
  const keys = [...counter.keys()].sort()

  let res = 0

  for (let i = 0; i < n; i++) {
    // 在当前位之后小于其的个数
    let smaller = 0
    // 这里应该用counter计数
    for (const key of keys) {
      if (key >= str[i]) break
      smaller += counter.get(key) ?? 0
    }

    res += factorial[n - 1 - i] * smaller
    counter.set(str[i], counter.get(str[i])! - 1)
  }

  return res
}

console.log(calRankByPerm(34152, 5))
console.log(calRankByPerm(1342, 4))

// 逆康托展开
// 第rank大排列，rank从1开始
// 瓶颈在删除的O(n) 使用sortedList可以O(nlogn)
function calPermByRank(rank: number, n: number): string {
  rank--
  const available = Array.from<unknown, number>({ length: n }, (_, i) => i + 1)
  const sb: number[] = []

  for (let i = n; i >= 1; i--) {
    const [div, mod] = [~~(rank / factorial[i - 1]), rank % factorial[i - 1]]
    rank = mod
    sb.push(available[div])
    available.splice(div, 1)
  }

  return sb.join('')
}

console.log(calPermByRank(62, 5))
console.log(calPermByRank(1, 5))
