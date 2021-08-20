/**
 * @param {string} s
 * @param {number} k
 * @return {string}
 */
const reverseStr = function (s: string, k: number): string {
  const len = s.length
  let resArr = s.split('')

  for (let i = 0; i < len; i += 2 * k) {
    let l = i - 1,
      r = Math.min(len, i + k)
    // 反转每一段
    while (++l < --r) [resArr[l], resArr[r]] = [resArr[r], resArr[l]]
  }
  return resArr.join('')
}

// 每 2k 个字符反转前 k 个字符。
console.log(reverseStr('abcdefg', 2))
