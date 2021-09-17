// 不使用 * 运算符， 实现两个正整数的相乘
function multiply(A: number, B: number): number {
  let res = 0
  let cur = 0
  while (A) {
    if (A & 1) res += B << cur
    A = A >> 1
    cur++
  }

  return res
}
