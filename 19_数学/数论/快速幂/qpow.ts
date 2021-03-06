/**
 * @returns a^n % mod
 */
function qpow(a: number, n: number, mod: number): number {
  if (n === 0) return 1 % mod

  let [_a, _n, _mod] = [BigInt(a), BigInt(n), BigInt(mod)]
  let res = 1n

  while (_n) {
    if (_n & 1n) {
      res *= _a
      res %= _mod
    }

    // 此处可能超出2^53-1 需要大数 (1e9-7*(1e9-7已经超出))
    _a *= _a
    _a %= _mod
    _n >>= 1n
  }

  return Number(res)
}

export { qpow }

// 如果a*a大于2^53-1则会丢失精度
// 注意js的快速幂可能需要大数
