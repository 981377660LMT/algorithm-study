/**
 * @param {number} num
 * @return {string}
 * 十进制小数是如何转换为二进制的
 * 小数是乘2，取整数，拿小数继续算，直到小数位是0

  0.8125x2=1.625 取整1,小数部分是0.625
  0.625x2=1.25 取整1,小数部分是0.25
  0.25x2=0.5 取整0,小数部分是0.5
  0.5x2=1.0 取整1,小数部分是0，结束。

  所以0.8125的二进制是0.1101

 */
const printBin = function (num: number): string {
  const res: string[] = ['0.']
  while (num) {
    num *= 2
    const integer = ~~num
    const decimal = num - integer
    res.push(integer.toString())
    num = decimal
    // 如果该数字无法精确地用32位以内的二进制表示，则打印“ERROR”。
    if (res.length >= 32) return 'ERROR'
  }

  return res.join('')
}

console.log(printBin(0.625))
// 输出："0.101"
console.log(printBin(0.1))
// 输出："ERROR"
// 提示：0.1无法被二进制准确表示
