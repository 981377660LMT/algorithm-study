/**
 * @param {string} version1
 * @param {string} version2
 * @return {number}
 * 左到右的顺序依次比较它们的修订号。比较修订号时，只需比较 忽略任何前导零后的整数值 。
 */
var compareVersion = function (version1: string, version2: string): number {
  const v1 = version1.split('.')
  const v2 = version2.split('.')
  const len = Math.max(v1.length, v2.length)
  for (let i = 0; i < len; i++) {
    const v1Version = parseInt(v1[i]) || 0
    const v2Version = parseInt(v2[i]) || 0
    if (v1Version === v2Version) continue
    return v1Version > v2Version ? 1 : -1
  }

  return 0
}

console.log(compareVersion('1.01', '1.001'))
// @ts-ignore
console.log(parseInt(undefined) || 0)
