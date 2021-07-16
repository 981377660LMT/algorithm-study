// 给定一个 32 位有符号整数，将整数中的数字进行反转。要求如下：

// 只翻转数字，符号位不进行翻转。
// 假设我们的环境只能存储 32 位有符号整数，其数值范围是 [-2^{31}, 2^{31}-1]。如果反转后的整数溢出，则返回 0。
// 不能借助JS原生的 reverse 函数
const reverseNumber = (n: number): number => {
  const metaInfo: { prefix: number; abs: number } = {
    prefix: n < 0 ? -1 : 1,
    abs: Math.abs(n),
  }

  const reversedNumber = metaInfo.abs
    .toString()
    .split('')
    .map(str => parseInt(str, 10))
    .reduce((pre, cur, index) => pre + cur * Math.pow(10, index), 0)

  if (n > Math.pow(2, 21) - 1 || n < Math.pow(-2, 31) + 1) return 0

  return metaInfo.prefix * reversedNumber
}

console.log(reverseNumber(-123))

export {}
