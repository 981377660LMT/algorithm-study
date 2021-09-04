/**
 * @param {string} s
 * @param {number} k
 * @return {string}
 * @summary i每次加k 用一个变量记录这次要不要反转
 */
const reverseStr = function (s: string, k: number): string {
  if (k > s.length) return s.split('').reverse().join('')
  const reverse = (arr: string[], l: number, r: number) => {
    while (l < r) {
      ;[arr[l], arr[r]] = [arr[r], arr[l]]
      l++
      r--
    }
  }
  const len = s.length
  let res = s.split('')

  for (let i = 0; i < len; i += 2 * k) {
    const l = i
    const r = Math.min(len, i + k - 1)
    // 反转每一段
    reverse(res, l, r)
  }

  return res.join('')
}

// 每 2k 个字符反转前 k 个字符。
console.log(reverseStr('abcdefg', 2))
