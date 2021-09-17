// 找出两个数字a和b中最大的那一个。不得使用if-else或其他比较运算符。
// 不超过32位时，正数有符号右移31位的时候值为0（正数的符号位为0），当负数有符号右移31位的时候值为-1（负数的符号位为1）
function maximum1(a: number, b: number): number {
  console.log(a, b)
  let diff = a - b
  const MOD = 10e9 + 7
  console.log(diff, diff % MOD, 777)
  const isMinus = (x: number) => !!(x >> 63)
  return isMinus(diff % MOD) ? b : a
}
function maximum2(a: number, b: number): number {
  return (Math.abs(b - a) + a + b) / 2
}
console.log(maximum2(2 ** 31 - 1, -1 * 2 ** 31))
// console.log(112 >> 31)
// console.log((2 ** 31 + 2 ** 31 + 1) >> 31)
// console.log(-13 >> 31)
