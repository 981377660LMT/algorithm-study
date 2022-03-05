/**
 * @param {string} num1
 * @param {string} num2
 * @return {string}
 */
function add(num1: string, num2: string): string {
  let i = num1.length - 1
  let j = num2.length - 1
  let carry = 0
  let res = ''

  while (carry || i >= 0 || j >= 0) {
    const d1 = i < 0 ? 0 : Number(num1[i])
    const d2 = j < 0 ? 0 : Number(num2[j])
    const sum = d1 + d2 + carry
    res = `${sum % 10}${res}`
    carry = ~~(sum / 10)
    i--
    j--
  }

  return res
}

if (require.main === module) {
  console.log(add('999999999999999999', '1'))
}

export { add }
