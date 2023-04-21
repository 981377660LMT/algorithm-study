// 两个数相乘取模 mod, 不使用BigInt
// O(logMod), 速度不如BigInt
function modMul(a: number, b: number, mod: number): number {
  let res = 0
  while (b) {
    if (b & 1) {
      res = (res + a) % mod
    }
    a = (a + a) % mod
    b >>>= 1
  }
  return res
}

export {}

if (require.main === module) {
  const a = 1e9 + 1
  const b = 1e9 + 6
  const MOD = 1e9 + 7
  console.log(modMul(a, b, MOD))
  console.log((a * b) % MOD)

  const n = 1e6
  const as1 = Array(n)
  const bs1 = Array(n)
  const as2 = Array(n)
  const bs2 = Array(n)
  for (let i = 0; i < n; i++) {
    as1[i] = Math.floor(Math.random() * MOD)
    bs1[i] = Math.floor(Math.random() * MOD)
    as2[i] = BigInt(as1[i])
    bs2[i] = BigInt(bs1[i])
  }

  const bigMod = BigInt(MOD)
  console.time('useBigIntToMul')
  for (let i = 0; i < n; i++) {
    const cur = BigInt(as2[i] * bs2[i]) % bigMod
  }
  console.timeEnd('useBigIntToMul')

  console.time('useModMul')
  for (let i = 0; i < n; i++) {
    const cur = modMul(as1[i], bs1[i], MOD)
  }
  console.timeEnd('useModMul')
}

// !BigInt还是比快速幂的乘法快的
