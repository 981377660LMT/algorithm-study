/* eslint-disable no-inner-declarations */
/* eslint-disable no-else-return */
/* eslint-disable newline-per-chained-call */

import { distSum } from '../../前缀与差分/template/distSum/distSum'

/**
 * 从小到大遍历`[min,max]`闭区间内的回文数.返回 true 可提前终止遍历.
 * @link https://github.com/EndlessCheng/codeforces-go
 */
function enumeratePalindrome(
  min: number,
  max: number,
  f: (palindrome: number) => boolean | void
): void {
  if (min > max) return
  const minLength = String(min).length
  const startBase = 10 ** ((minLength - 1) >>> 1)
  for (let base = startBase; ; base *= 10) {
    // 生成奇数长度回文数，例如 base = 10，生成的范围是 101 ~ 999
    for (let i = base; i < base * 10; i++) {
      let x = i
      for (let t = Math.floor(i / 10); t > 0; t = Math.floor(t / 10)) {
        x = x * 10 + (t % 10)
      }
      if (x > max) return
      if (x >= min) {
        if (f(x)) return
      }
    }

    // 生成偶数长度回文数，例如 base = 10，生成的范围是 1001 ~ 9999
    for (let i = base; i < base * 10; i++) {
      let x = i
      for (let t = i; t > 0; t = Math.floor(t / 10)) {
        x = x * 10 + (t % 10)
      }
      if (x > max) return
      if (x >= min) {
        if (f(x)) return
      }
    }
  }
}

/**
 * 遍历长度在 `[minLength, maxLength]` 之间的回文数字字符串.
 * @param maxLength maxLength <= 12.
 */
function enumeratePalindromeByLength(
  minLength: number,
  maxLength: number,
  f: (palindrome: string) => boolean | void,
  options: {
    reverse?: boolean
  } = {}
): void {
  if (minLength > maxLength) return
  const { reverse = false } = options

  if (reverse) {
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
 * 遍历长度在 `[minLength, maxLength]` 之间的回文数字字符串.
 * @param maxLength maxLength <= 12.
 */
function* generatePalindrome(
  minLength: number,
  maxLength: number,
  reversed = false
): Generator<string> {
  if (minLength > maxLength) return

  if (reversed) {
    for (let len = maxLength; len >= minLength; len--) {
      const start = 10 ** ((len - 1) >> 1)
      const end = start * 10 - 1
      for (let half = end; half >= start; half--) {
        if (len & 1) {
          yield `${half}${String(half).slice(0, -1).split('').reverse().join('')}`
        } else yield `${half}${String(half).split('').reverse().join('')}`
      }
    }
  } else {
    for (let len = minLength; len <= maxLength; len++) {
      const start = 10 ** ((len - 1) >> 1)
      const end = start * 10 - 1
      for (let half = start; half <= end; half++) {
        if (len & 1) {
          yield `${half}${String(half).slice(0, -1).split('').reverse().join('')}`
        } else yield `${half}${String(half).split('').reverse().join('')}`
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

/** 返回比x大的下一个回文数. */
function nextPalindrome(x: string): string {
  if (x === '9'.repeat(x.length)) return `1${'0'.repeat(x.length - 1)}1`
  if (x.length & 1) {
    const half = String(Number(x.slice(0, (x.length >>> 1) + 1)) + 1)
    return `${half}${half.slice(0, -1).split('').reverse().join('')}`
  } else {
    const half = String(Number(x.slice(0, x.length >>> 1)) + 1)
    return `${half}${half.split('').reverse().join('')}`
  }
}

export {
  enumeratePalindrome,
  enumeratePalindromeByLength,
  generatePalindrome,
  getPalindromeByHalf,
  countPalindrome,
  getKthPalindrome,
  nextPalindrome
}

if (require.main === module) {
  // 100151. 使数组成为等数数组的最小代价
  // https://leetcode.cn/problems/minimum-cost-to-make-array-equalindromic/description/
  const palindromes: number[] = []
  enumeratePalindrome(1, 10 ** 9, p => {
    palindromes.push(p)
  })

  function minimumCost(nums: number[]): number {
    nums.sort((a, b) => a - b)
    const D = distSum(nums)
    let res = Infinity
    // !注意不要用 Math.min(...palindromes.map(D))，会超时
    for (let i = 0; i < palindromes.length; i++) {
      res = Math.min(res, D(palindromes[i]))
    }
    return res
  }
}
