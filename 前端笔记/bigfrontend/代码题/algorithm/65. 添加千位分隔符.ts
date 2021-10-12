/**
 * @param {number} num
 * @return {string}
 */
function addComma(num: number): string {
  const [integer, decimal] = num.toString().split('.')

  const sb = integer.toString().split('')
  for (let i = sb.length - 4; i >= 0; i -= 3) {
    sb[i] = sb[i] + ','
  }

  return sb.join('') + (decimal ? `.${decimal}` : '')
}

console.log(addComma(12345678.12345))
