export {}

/**
 * 将字符串中的指定子串替换为另一个子串.
 *
 * @param str 要替换的字符串.
 * @param oldStr 要替换的子串.
 * @param newStr 替换后的子串.
 * @param limit 替换的次数限制.
 * @returns 替换后的字符串.
 */
function myReplace(str: string, oldStr: string, newStr: string, limit?: number): string {
  if (limit === 0) return str
  if (limit === undefined) limit = Infinity
  if (oldStr === '') return str
  const res: string[] = []
  let ptr = 0
  for (let i = 0; i < limit; i++) {
    const pos = str.indexOf(oldStr, ptr)
    if (pos === -1) {
      break
    }
    res.push(str.slice(ptr, pos))
    res.push(newStr)
    ptr = pos + oldStr.length
  }
  res.push(str.slice(ptr))
  return res.join('')
}

console.log(myReplace('a,b,c', ',', '!', 0)) // a,b,c
console.log(myReplace('a,b,c', ',', '!', 1)) // a!b,c
console.log(myReplace('a,b,c', ',', '!', 2)) // a!b!c
console.log(myReplace('a,b,c', ',', '!', 3)) // a!b!c
