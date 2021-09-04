/**
 * @param {string} s
 * @param {number} k
 * @return {string}
 * 第一个分组包含的字符个数必须小于等于 K，
 * 但至少要包含 1 个字符。两个分组之间需要用 '-'（破折号）隔开
 * @summary
 * 使用'-'分割相当于在那一项数组元素加一个'-'
 * 倒序改变数组的项
 */
var licenseKeyFormatting = function (s: string, k: number): string {
  const chars = s.replace(/-/g, '').toUpperCase().split('')
  // 注意这个边界 +k 就到了最后一个字符 不需要加'-'
  for (let i = chars.length - 1 - k; i >= 0; i -= k) {
    chars[i] += '-'
  }

  return chars.join('')
}

console.log(licenseKeyFormatting('5F3Z-2e-9-w', 4))

export default 1
// 输出："5F3Z-2E9W"
// 解释：字符串 S 被分成了两个部分，每部分 4 个字符；
//      注意，两个额外的破折号需要删掉。
