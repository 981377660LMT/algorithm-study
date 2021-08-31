/**
 * @param {string} num1  num1 和num2 的长度都小于 5100
 * @param {string} num2
 * @return {string}
 * 你不能使用任何內建 BigInteger 库， 也不能直接将输入的字符串转换为整数形式
 */
const addStrings = function (num1: string, num2: string): string {
  let i = num1.length - 1
  let j = num2.length - 1
  let carry = 0
  let res = ''

  while (i >= 0 || j >= 0) {
    const d1 = i < 0 ? 0 : parseInt(num1[i])
    const d2 = j < 0 ? 0 : parseInt(num2[j])
    const sum = d1 + d2 + carry
    res = `${sum % 10}${res}`
    carry = ~~(sum / 10)
    i--
    j--
  }

  carry && (res = `1${res}`)
  return res
}

console.log(addStrings('11', '123'))

export {}
