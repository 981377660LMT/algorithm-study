/**
 * @param {string} p  p 的数据范围是 10^5 ，因此暴力找出所有子串就需要 10^10 次操作了
 * @return {number}
 * 找出 s 中有多少个唯一的 p 的非空子串
 * s 看作是“abcdefghijklmnopqrstuvwxyz”的无限环绕字符串，所以 s 看起来是这样的："...zabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcd....".
 */
const findSubstringInWraproundString = function (p: string): number {
  // p = '$' + p
  let pre = 1
  // 以某个字母结尾的子串个数（bcabc 则取6=1+2+3 因为 abc包含了bc的情况,最长的连续子串一定是包含了比它短的连续子串),防止重复计算
  const mapper = new Map<string, number>([[p[0], 1]])

  for (let i = 1; i < p.length; i++) {
    if ([1, -25].includes(p[i].codePointAt(0)! - p[i - 1].codePointAt(0)!)) pre++
    else pre = 1
    mapper.set(p[i], Math.max(pre, mapper.get(p[i]) || 1))
  }

  console.log(mapper)
  return Array.from(mapper.values()).reduce((pre, cur) => pre + cur, 0)
}

console.log(findSubstringInWraproundString('zab'))
console.log(findSubstringInWraproundString('cac'))
console.log(findSubstringInWraproundString('a'))

export {}
