/**
 * @param {string} pattern  1 <= len(pattern) <= 1000
 * @param {string} value  0 <= len(value) <= 1000
 * @return {boolean}
 * @description
 * pattern只包含字母"a"和"b"，value仅包含小写字母。
 */
// const patternMatching = function (pattern: string, value: string): boolean {}

const patternMatching = function (pattern: string, value: string): boolean {
  if (pattern.length <= 1) return true
  if (!value.length) return false
  let regexpPattern = ''
  const [regA, regB] = pattern[0] === 'a' ? ['\\1', '\\2'] : ['\\2', '\\1']
  regexpPattern = pattern
    .replace('a', '(.*?)') // 只会替换一个
    .replace('b', '(.*?)')
    .replace(/a/g, regA) // 替换全部
    .replace(/b/g, regB)
  regexpPattern = '^' + regexpPattern + '$'
  const regexp = new RegExp(regexpPattern, 'g')
  return regexp.test(value)
}

console.log(patternMatching('abba', 'dogcatcatdog')) // true
console.log(patternMatching('abba', 'dogcatcatfish')) // false
