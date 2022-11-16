// js 怎么写记忆化最快
// !Int32Array 静态数组存储结果 不过要注意Int32Array最多存2e9

const MOD = 1e9 + 7

const gcd = (a: number, b: number): number => {
  if (Number.isNaN(a) || Number.isNaN(b)) return NaN
  return b === 0 ? a : gcd(b, a % b)
}

function distinctSequences(n: number): number {
  const dfs = (index: number, pre1: number, pre2: number): number => {
    if (index === n) return 1

    const hash = pre1 * 7 + pre2
    if (cache[index][hash] !== -1) return cache[index][hash]

    let res = 0
    for (let cur = 1; cur <= 6; cur++) {
      if (pre1 === cur || pre2 === cur) continue
      if (pre2 === 0 || gcd(cur, pre2) === 1) {
        res += dfs(index + 1, pre2, cur)
        res %= MOD
      }
    }

    cache[index][hash] = res
    return res
  }

  const cache = Array.from({ length: n }, () => new Int32Array(100).fill(-1))
  return dfs(0, 0, 0)
}

export {}
