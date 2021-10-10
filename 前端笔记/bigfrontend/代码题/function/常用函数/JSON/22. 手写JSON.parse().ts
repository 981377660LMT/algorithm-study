/**
 * @param {string} str
 * @return {object | Array | string | number | boolean | null}
 */
function parse(str: string): object | Array<any> | string | number | boolean | null {
  // your code here
  if (str === '') throw Error() // 空输入会报错

  if (str[0] === "'") throw Error()

  if (str === 'null') {
    return null
  }
  if (str === '{}') {
    return {}
  }
  if (str === '[]') {
    return []
  }
  if (str === 'true') {
    return true
  }
  if (str === 'false') {
    return false
  }
  if (str[0] === '"') {
    return str.slice(1, -1)
  }

  if (+str === +str) {
    return Number(str)
  }

  if (str[0] === '[') {
    return str
      .slice(1, -1)
      .split(',')
      .map(value => parse(value))
  }

  if (str[0] === '{') {
    return str
      .slice(1, -1)
      .split(',')
      .reduce<Record<PropertyKey, any>>((acc, item) => {
        const index = item.indexOf(':')
        const key = item.slice(0, index)
        const value = item.slice(index + 1)
        acc[parse(key) as PropertyKey] = parse(value)
        return acc
      }, {})
  }

  throw new Error('invalid input')
}

// // 晕了
