/**
 * @param {string} num1
 * @param {string} num2
 * @return {string}
 * 所有输入都是非负整数字符串，结果也保证是非负的。
 */
function subtract(num1: string, num2: string): string {
  let i = num1.length - 1
  let j = num2.length - 1
  let carry = 0
  let res = ''

  while (i >= 0 || j >= 0) {
    const d1 = i < 0 ? 0 : parseInt(num1[i])
    const d2 = j < 0 ? 0 : parseInt(num2[j])
    const diff = d1 - d2 - carry
    res = `${(diff + 10) % 10}${res}`
    carry = diff < 0 ? 1 : 0
    i--
    j--
  }

  // 保证结果非负 所以不用处理最后的carry
  // carry && (res = `1${res}`)
  return res.replace(/^0+/g, '') || '0'
}

if (require.main === module) {
  console.log(subtract('1000000000000000000000', '999999999999999999999'))
}

export { subtract }
