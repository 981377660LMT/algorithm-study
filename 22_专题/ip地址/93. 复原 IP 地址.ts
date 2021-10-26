/**
 * @param {string} s
 * @return {string[]}
 */
function restoreIpAddresses(s: string): string[] {
  if (s.length <= 3) return []
  if (s.length >= 13) return []

  const res: string[][] = []

  // 有效 IP 地址 正好由四个整数（每个整数位于 0 到 255 之间组成，且不能含有前导 0），整数之间用 '.' 分隔。
  const isValidString = (str: string) => {
    if (parseInt(str, 10) > 255) return false
    if (str.length >= 2 && str[0] === '0') return false
    if (str.length === 0 || str.length > 3) return false
    return true
  }

  const bt = (path: string[], subString: string) => {
    // 1. 终点条件
    if (path.length === 3) {
      if (isValidString(subString)) {
        return res.push([...path, subString])
      }
    }

    for (let i = 1; i <= 3; i++) {
      const sub = subString.slice(0, i)
      // 2.排除不合的
      if (isValidString(sub)) {
        // 3.继续bt
        path.push(sub)
        bt(path, subString.slice(i))
        path.pop()
      }
    }
  }
  bt([], s)

  return res.map(strArr => strArr.join('.'))
}
console.log(restoreIpAddresses('101023'))
