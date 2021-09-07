/**
 * @param {number} numerator
 * @param {number} denominator
 * @return {string}
 * 以 字符串形式返回小数 。
   如果小数部分为循环小数，则将循环的部分括在括号内。
   这种题有几种情况

   1.正负号问题
   2.加小数点的情况, 比如 8/ 2 不需要加小数点
   3.小数部分,如何判断是否开始循环了

   1.先判断结果的正负
   2.直接相除, 通过余数,看能否整除
   3.开始循环的时候, 说明之前已经出现过这个余数, 我们只要记录前面出现余数的位置,插入括号即可!
 */
const fractionToDecimal = function (numerator: number, denominator: number): string {
  if (numerator === 0) return '0'
  const res: string[] = []

  // 处理符号问题
  const isNegative = numerator > 0 !== denominator > 0
  isNegative && res.push('-')
  ;[numerator, denominator] = [Math.abs(numerator), Math.abs(denominator)]

  // 处理小数问题  这里不能用 ~~ 会出现符号
  let [div, mod] = [Math.floor(numerator / denominator), numerator % denominator]
  console.log(div, mod)
  res.push(div.toString())
  if (mod === 0) return res.join('')
  res.push('.')

  // 处理余数:模拟除法的过程
  const modStartIndex = new Map<number, number>([[mod, res.length]])
  while (mod) {
    mod *= 10
    ;[div, mod] = [Math.floor(mod / denominator), mod % denominator]
    res.push(div.toString())
    if (modStartIndex.has(mod)) {
      res.splice(modStartIndex.get(mod)!, 0, '(')
      res.push(')')
      break
    }
    modStartIndex.set(mod, res.length)
  }

  return res.join('')
}

// console.log(fractionToDecimal(1, 2)) // "0.5"
// console.log(fractionToDecimal(4, 333)) // "0.(012)"
console.log(fractionToDecimal(-2147483648, 1)) // "0.(012)"

export default 1

// "在"1号索引上插入
const a = [1, 2, 3] as any[]
a.splice(1, 0, '(')
console.log(a)
