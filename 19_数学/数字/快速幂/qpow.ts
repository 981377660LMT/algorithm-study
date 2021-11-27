const MOD = 1e9 + 7
function qpow(a: number, n: number): number {
  let res = 1

  while (n) {
    if (n & 1) {
      res *= a
      // res %= MOD
    }

    a *= a
    // a %= MOD
    n >>>= 1
  }

  return res
}

export { qpow }

// 如果a*a大于2^53-1则会丢失精度
// 注意js的快速幂可能需要大数
