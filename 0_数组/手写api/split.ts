export {}

/**
 * 将字符串分割成数组.
 *
 * @param str 要分割的字符串.
 * @param separator 分隔符.
 * @param limit 返回的数组长度限制.
 * @returns 分割后的数组.
 *
 * @example
 * ```ts
 * mySplit('a,b,c', ',', 0) // []
 * mySplit('a,b,c', ',', 1) // ['a']
 * mySplit('a,b,c', ',', 2) // ['a', 'b']
 * mySplit('a,b,c', ',', 3) // ['a', 'b', 'c']
 * ```
 */
function mySplit(str: string, separator: string, limit?: number): string[] {
  if (limit === 0) return []
  if (limit === undefined) limit = Infinity

  const res: string[] = []

  if (separator === '') {
    for (let i = 0; i < Math.min(str.length, limit); i++) {
      res.push(str[i])
    }
    return res
  }

  let ptr = 0
  let sepIndex = 0
  while (res.length < limit) {
    sepIndex = str.indexOf(separator, ptr)
    if (sepIndex === -1) {
      break
    }
    res.push(str.slice(ptr, sepIndex))
    ptr = sepIndex + separator.length
  }
  if (res.length < limit) {
    res.push(str.slice(ptr))
  }
  return res
}
