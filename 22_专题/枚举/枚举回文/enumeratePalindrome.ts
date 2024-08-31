/* eslint-disable no-inner-declarations */
/* eslint-disable no-else-return */
/* eslint-disable newline-per-chained-call */

/**
 * 从小到大遍历`[min,max]`闭区间内的回文数.返回 true 可提前终止遍历.
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
      for (let t = (i / 10) | 0; t > 0; t = (t / 10) | 0) {
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
      for (let t = i; t > 0; t = (t / 10) | 0) {
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
 * 遍历数位长度在 `[minLength, maxLength]` 之间的回文数.
 * @param maxLength maxLength <= 14.
 */
function enumeratePalindromeByLength(
  minLength: number,
  maxLength: number,
  f: (palindrome: number) => boolean | void
): void {
  if (minLength > maxLength) return
  const min = 10 ** (minLength - 1)
  const max = 10 ** maxLength - 1
  const startBase = 10 ** ((minLength - 1) >>> 1)
  for (let base = startBase; ; base *= 10) {
    for (let i = base; i < base * 10; i++) {
      let x = i
      for (let t = (i / 10) | 0; t > 0; t = (t / 10) | 0) {
        x = x * 10 + (t % 10)
      }
      if (x > max) return
      if (x >= min) {
        if (f(x)) return
      }
    }

    for (let i = base; i < base * 10; i++) {
      let x = i
      for (let t = i; t > 0; t = (t / 10) | 0) {
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
  getPalindromeByHalf,
  countPalindrome,
  getKthPalindrome,
  nextPalindrome
}

if (require.main === module) {
  // 3272. 统计好整数的数目
  // https://leetcode.cn/problems/find-the-count-of-good-integers/description/
  function countGoodIntegers(n: number, k: number): number {
    const FAC: number[] = [1]
    for (let i = 1; i <= 15; i++) FAC.push(FAC[i - 1] * i)
    const counter = Array(10).fill(0)

    let res = 0
    const visited: Set<string> = new Set()
    enumeratePalindromeByLength(n, n, palindrome => {
      if (palindrome % k === 0) {
        const key = String(palindrome).split('').sort().join('')
        if (!visited.has(key)) {
          visited.add(key)
          res += calc(palindrome)
        }
      }
    })
    return res

    /**
     * 字符串重新排列不含前导0的数字的个数.
     */
    function calc(palindrome: number): number {
      counter.fill(0)
      let cur = palindrome
      let m = 0
      while (cur) {
        counter[cur % 10]++
        cur = (cur / 10) | 0
        m++
      }

      let res = (m - counter[0]) * FAC[m - 1]
      for (let i = 0; i < 10; i++) {
        res /= FAC[counter[i]]
      }
      return res
    }
  }
}
