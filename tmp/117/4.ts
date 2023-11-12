export {}

const MOD = 1e9 + 7
function stringCount(n: number): number {
  if (n <= 3) {
    return 0
  }

  const memo = new Int32Array(n * 3 * 3 * 3).fill(-1)
  const dfs = (index: number, a: number, b: number, c: number): number => {
    if (index === n) {
      return a >= 1 && b >= 2 && c >= 1 ? 1 : 0
    }

    const hash = index * 27 + a * 9 + b * 3 + c
    if (memo[hash] !== -1) return memo[hash]

    let res = dfs(index + 1, a, b, c) * 23
    res += dfs(index + 1, Math.min(a + 1, 1), b, c)
    res += dfs(index + 1, a, Math.min(b + 1, 2), c)
    res += dfs(index + 1, a, b, Math.min(c + 1, 1))
    res %= MOD

    memo[hash] = res
    return res
  }

  return dfs(0, 0, 0, 0)
}

let n: number
let k: number
const
