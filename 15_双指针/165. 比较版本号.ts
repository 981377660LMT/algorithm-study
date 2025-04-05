/**
 * @param {string} version1
 * @param {string} version2
 * @return {number}
 * 左到右的顺序依次比较它们的修订号。比较修订号时，只需比较 忽略任何前导零后的整数值 。
 * 比较版本号时，请按 从左到右的顺序 依次比较它们的修订号。如果其中一个版本字符串的修订号较少，则将缺失的修订号视为 0。
 */
const compareVersion = function (version1: string, version2: string): number {
  const v1 = version1.split('.')
  const v2 = version2.split('.')
  const len = Math.max(v1.length, v2.length)
  for (let i = 0; i < len; i++) {
    const n1 = i < v1.length ? Number(v1[i]) : 0
    const n2 = i < v2.length ? Number(v2[i]) : 0
    if (n1 !== n2) {
      return n1 > n2 ? 1 : -1
    }
  }
  return 0
}

console.log(compareVersion('1.01', '1.001'))
console.log(compareVersion('7.5.2.4', '7.5.3'))
// @ts-ignore
console.log(parseInt(undefined) || 0)

export {}
