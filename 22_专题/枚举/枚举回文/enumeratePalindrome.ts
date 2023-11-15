/* eslint-disable newline-per-chained-call */

/**
 * 遍历长度在 `[minLength, maxLength]` 之间的回文数字字符串.
 * @param maxLength maxLength <= 12.
 */
function enumeratePalindrome(minLength: number, maxLength: number, f: (palindrome: string) => boolean | void, reversed = false): void {
  if (minLength > maxLength) return

  if (reversed) {
    for (let len = maxLength; len >= minLength; len--) {
      const start = 10 ** ((len - 1) >> 1)
      const end = start * 10 - 1
      for (let half = end; half >= start; half--) {
        if (len & 1) {
          if (f(`${half}${String(half).slice(0, -1).split('').reverse().join('')}`)) return
        } else if (f(`${half}${String(half).split('').reverse().join('')}`)) return
      }
    }
  } else {
    for (let len = minLength; len <= maxLength; len++) {
      const start = 10 ** ((len - 1) >> 1)
      const end = start * 10 - 1
      for (let half = start; half <= end; half++) {
        if (len & 1) {
          if (f(`${half}${String(half).slice(0, -1).split('').reverse().join('')}`)) return
        } else if (f(`${half}${String(half).split('').reverse().join('')}`)) return
      }
    }
  }
}

/**
 * 给定回文的一半,返回偶数长度/奇数长度的回文字符串.
 */
function getPalindromeByHalf(half: string | number, even = true): string {
  if (even) return `${half}${String(half).split('').reverse().join('')}`
  return `${half}${String(half).slice(0, -1).split('').reverse().join('')}`
}

/**
 * 返回长度为length的回文数个数.
 */
function countPalindrome(length: number): number {
  if (length <= 0) return 0
  const start = 10 ** ((length - 1) >> 1)
  return start * 10 - 1 - start + 1
}

/**
 * 返回长度为length的第k个回文数,k>=1.
 */
function getKthPalindrome(length: number, k: number): string | undefined {
  if (length <= 0) return undefined
  const start = 10 ** ((length - 1) >> 1)
  const count = start * 10 - 1 - start + 1
  if (k > count) return undefined
  const half = start + k - 1
  if (length & 1) return `${half}${String(half).slice(0, -1).split('').reverse().join('')}`
  return `${half}${String(half).split('').reverse().join('')}`
}

export { enumeratePalindrome, getPalindromeByHalf, countPalindrome, getKthPalindrome }

if (require.main === module) {
  let count = 0
  enumeratePalindrome(1, 13, p => {
    count++
  })
  console.log(count) // 1999998
}
