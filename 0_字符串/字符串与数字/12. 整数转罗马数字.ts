// 整数转罗马数字

const ROMAN = ['M', 'CM', 'D', 'CD', 'C', 'XC', 'L', 'XL', 'X', 'IX', 'V', 'IV', 'I']
const ARAB = [1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1]

/**
 * @param {number} num
 * @return {string}
 * @description 类似于n进制转换
 * 除 模 除 模 ...
 */
function intToRoman(num: number): string {
  const res: string[] = []
  let radixIndex = 0
  while (num) {
    const [div, mod] = [~~(num / ARAB[radixIndex]), num % ARAB[radixIndex]]
    res.push(ROMAN[radixIndex].repeat(div))
    num = mod
    radixIndex++
  }
  return res.join('')
}

console.log(intToRoman(9))

export default 1
