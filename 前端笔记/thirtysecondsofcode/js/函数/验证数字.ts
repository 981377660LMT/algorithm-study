import assert from 'assert'

/**
 * @description javascript里判断一个字符串是不是由纯数字组成需要两步：
 * parseFloat 的结果不为NaN、原字符串不为正负 Infinity
 * @waring '1e6'会被返回true
 */
function isNumerical(str: string): boolean {
  if (typeof str !== 'string') return false
  const num = parseFloat(str)
  // 注意不要用Number.isFinite(str)/Number.isNaN(str) 因为会先把字符串转成数字 从而所有字符串都被认为false
  return !isNaN(num) && isFinite(Number(str))
}

/**
 * @description 判断是否只包含数字
 * 类似科学记数法，比如1e3，parseInt之类的是可以解析的，正则这个会得到false
 */
function isDigit(str: string): boolean {
  return /^\d+$/.test(str)
}

function isPositiveInteger(num: number): boolean {
  if (typeof num !== 'number') return false
  return num > 0 && Number.isInteger(num)
}

if (require.main === module) {
  console.log(isNumerical('10')) // true
  console.log(isNumerical('11*'))
  console.log(isNumerical(''))

  assert.strictEqual(isNumerical(''), false)
  assert.strictEqual(isNumerical('1e6'), true)
  assert.strictEqual(isNumerical('12.9'), true)
  assert.strictEqual(isNumerical('Infinity'), false)
  assert.strictEqual(isNumerical('12*'), false)
  assert.strictEqual(isDigit('1e5'), false)
}

export { isNumerical, isDigit, isPositiveInteger }
