/**
 * @param {string[]} digits
 * @param {number} n 1 <= N <= 10^9
 * @return {number}
 * 我们有一组排序的数字 D，它是  {'1','2','3','4','5','6','7','8','9'} 的非空子集。
 * 返回可以用 D 中的数字写出的小于或等于 N 的正整数的数目。
 * @description 必须用数学方法 否则超时
 */
const atMostNGivenDigitSet = function (digits: string[], n: number): number {
  const nums = digits.map(Number)
  const len = n.toString().length
  let res = 0

  // 小于len位的
  for (let i = 1; i < len; i++) {
    res += nums.length ** i
  }

  // 等于len位时 n的高位向低位对比
  for (let i = 0; i < len; i++) {
    const upper = Number(n.toString()[i])
    // 在此位能取的数
    const lessThan = nums.filter(v => v < upper).length
    res += lessThan * digits.length ** (len - i - 1)
    console.log(res, lessThan)

    // 如果存在upper 则不必继续看了比n小的已经看完
    if (!nums.includes(upper)) break
    if (i === len - 1) res++ // 等于n
  }

  return res
}

console.log(atMostNGivenDigitSet(['1', '3', '5', '7'], 786))
console.log(atMostNGivenDigitSet(['1', '3', '5', '7'], 100))

export {}
