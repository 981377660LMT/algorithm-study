const N = 1e5
const MOD = BigInt(1e9 + 7)
const fac = Array<bigint>(N + 1).fill(1n)
const ifac = Array<bigint>(N + 1).fill(1n)

for (let i = 2; i <= N; i++) {
  fac[i] = (fac[i - 1] * BigInt(i)) % MOD
  ifac[i] = calInv(fac[i], MOD)
}

function comb(n: number, k: number): number {
  if (n < k) return 1
  return Number((((fac[n] * ifac[k]) % MOD) * ifac[n - k]) % MOD)
}

function calInv(a: bigint, mod: bigint) {
  return qpow(a, mod - 2n)
}

// 注意js精度丢失问题 需要用大数
function qpow(a: bigint, n: bigint) {
  if (n === 0n) return 1n % MOD

  let res = 1n

  while (n) {
    if (n & 1n) {
      res *= a
      res %= MOD
    }

    // 此处可能超出2^53-1 需要大数 (1e9-7*(1e9-7已经超出))
    a *= a
    a %= MOD
    n >>= 1n
  }

  return res
}

if (require.main === module) {
  console.log(fac, ifac)
  console.log(comb(10, 1) % Number(MOD))
  console.log(calInv(2n, MOD))
  console.log((2n * calInv(2n, MOD)) % MOD)
  console.log(qpow(2n, 1000000n))
  console.log(Number(BigInt(2 ** 53 + 1)))
  console.log(Number(BigInt(2 ** 53 - 1)))
  console.log(Number(BigInt(2 ** 100)))
}

export { comb }
