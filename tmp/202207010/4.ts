const LOWERCASE = 'abcdefghijklmnopqrstuvwxyz'
const UPPERCASE = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
const DIGITS = '0123456789'
const MOD = 1e9 + 7
const EPS = 1e-8
function factors(n: number): number[] {
  if (n <= 0) return []

  const small: number[] = []
  const big: number[] = []

  const upper = Math.floor(Math.sqrt(n))
  for (let f = 1; f <= upper; f++) {
    if (n % f === 0) {
      small.push(f)
      big.push(n / f)
    }
  }

  if (small[small.length - 1] === big[big.length - 1]) big.pop()

  return [...small, ...big.reverse()]
}

const F: number[][] = []
for (let i = 0; i <= 1e4; i++) {
  F.push(factors(i))
}

const cache = Array.from({ length: 1e4 + 1 }, () => new Int32Array(1e4 + 1).fill(-1))
function dfs(remain: number, cur: number): number {
  if (remain === 0) return 1
  if (cache[remain][cur] !== -1) return cache[remain][cur]
  let res = 0
  for (const next of F[cur]) {
    res += dfs(remain - 1, next)
    res %= MOD
  }
  cache[remain][cur] = res
  return res
}

function idealArrays(n: number, maxValue: number): number {
  let res = 0
  for (let i = 1; i <= maxValue; i++) {
    res += dfs(n - 1, i)
    res %= MOD
  }
  return res
}
