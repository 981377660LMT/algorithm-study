/* eslint-disable no-param-reassign */
/* eslint-disable no-shadow */

// 注意质因数分解/因数分解的上界不要写sqrt(n)，要写 f*f <=n .
// a%-b == a%b

import assert from 'assert'

/** 埃氏筛. */
class EratosthenesSieve {
  /**
   * 每个数的最小质因子.
   */
  readonly minPrime: Uint32Array
  private readonly _max: number

  constructor(max: number) {
    const minPrime = new Uint32Array(max + 1)
    for (let i = 0; i <= max; i++) minPrime[i] = i
    const upper = ~~Math.sqrt(max)
    for (let i = 2; i <= upper; i++) {
      if (minPrime[i] < i) continue
      for (let j = i * i; j <= max; j += i) {
        if (minPrime[j] === j) minPrime[j] = i
      }
    }
    this.minPrime = minPrime
    this._max = max
  }

  isPrime(n: number): boolean {
    return n >= 2 && this.minPrime[n] === n
  }

  getPrimeFactors(n: number): ReadonlyMap<number, number> {
    const f = this.minPrime
    const res = new Map<number, number>()
    while (n > 1) {
      const p = f[n]
      res.set(p, (res.get(p) || 0) + 1)
      n /= p
    }
    return res
  }

  getPrimes(n = this._max): readonly number[] {
    const res: number[] = []
    for (let i = 2; i <= n; i++) {
      if (i === this.minPrime[i]) res.push(i)
    }
    return res
  }
}

// O(n^0.5)
function isPrime(n: number): boolean {
  if (n < 2 || (n >= 4 && n % 2 === 0)) return false
  const upper = ~~Math.sqrt(n)
  for (let i = 2; i <= upper; i++) {
    if (n % i === 0) {
      return false
    }
  }
  return true
}

/**
 * @returns 返回 n 的所有质数因子，键为质数，值为因子的指数。
 * O(n^0.5)
 */
function getPrimeFactors(n: number): Map<number, number> {
  if (n <= 1) return new Map()
  const factors = new Map<number, number>()
  let count2 = 0
  while (n % 2 === 0) {
    n /= 2
    count2++
  }
  if (count2) factors.set(2, count2)
  // 跳过偶数
  for (let i = 3; i * i <= n; i += 2) {
    let count = 0
    while (n % i === 0) {
      n /= i
      count++
    }
    if (count) factors.set(i, count)
  }
  if (n > 1) factors.set(n, 1)
  return factors
}

/**
 * 区间质数个数.
 * [floor, ceiling]内的质数个数.
 * 1<=floor<=ceiling<=1e12,ceiling-floor<=5e5
 */
function countPrime(floor: number, ceiling: number): number {
  const isPrime = new Uint8Array(ceiling - floor + 1)
  for (let i = 0; i < isPrime.length; i++) isPrime[i] = 1
  isPrime[0] = +(floor !== 1)

  const last = ~~Math.sqrt(ceiling)
  for (let fac = 2; fac <= last; fac++) {
    let start = fac * Math.max(0 | (floor + fac - 1 / fac), 2) - floor
    while (start < isPrime.length) {
      isPrime[start] = 0
      start += fac
    }
  }

  let res = 0
  isPrime.forEach(v => {
    res += v
  })
  return res
}

/**
 * 区间筛/分段筛求 [floor,higher) 中的每个数是否为质数.
 * 1<=floor<=higher<=1e12,higher-floor<=5e5.
 */
function segmentedSieve(floor: number, higher: number): boolean[] {
  let root = 1
  while ((root + 1) * (root + 1) < higher) {
    root++
  }

  const isPrime = new Uint8Array(root + 1)
  for (let i = 0; i < isPrime.length; i++) isPrime[i] = 1
  isPrime[0] = 0
  isPrime[1] = 0

  const res: boolean[] = Array(higher - floor)
  for (let i = 0; i < res.length; i++) res[i] = true
  for (let i = 0; i < 2 - floor; i++) res[i] = false

  for (let i = 2; i <= root; i++) {
    if (isPrime[i]) {
      for (let j = i * i; j <= root; j += i) {
        isPrime[j] = 0
      }
      for (let j = Math.max(0 | ((floor + i - 1) / i), 2) * i; j < higher; j += i) {
        res[j - floor] = false
      }
    }
  }

  return res
}

/**
 * 返回 n 的所有因子.
 * @complexity O(n^0.5)
 */
function getFactors(n: number): number[] {
  if (n <= 0) return []
  const big: number[] = []
  const small: number[] = []
  const upper = Math.floor(Math.sqrt(n))
  for (let i = 1; i <= upper; i++) {
    if (n % i === 0) {
      small.push(i)
      if (i * i < n) big.push(n / i)
    }
  }
  for (let i = big.length - 1; i >= 0; i--) {
    small.push(big[i])
  }
  return small
}

/**
 * 空间复杂度为O(1)的枚举因子.枚举顺序为从小到大.
 */
function enumerateFactors(n: number, f: (factor: number) => void): void {
  if (n <= 0) return
  let i = 1
  const upper = Math.floor(Math.sqrt(n))
  for (; i <= upper; i++) {
    if (n % i === 0) f(i)
  }
  i--
  if (i * i === n) i--
  for (; i > 0; i--) {
    if (n % i === 0) f(n / i)
  }
}

/**
 * 返回区间 `[0, upper]` 内所有数的约数.
 * @param upper 上界.
 */
function getFactorsOfAll(upper: number): number[][] {
  const res: number[][] = Array(upper + 1)
  for (let i = 0; i <= upper; i++) res[i] = []
  for (let i = 1; i <= upper; i++) {
    for (let j = i; j <= upper; j += i) res[j].push(i)
  }
  return res
}

/**
 * 返回约数个数.
 * @param primeFactors 质因子分解.如果分解为空，则返回 1.
 */
function countFactors(nOrprimeFactors: number | ReadonlyMap<number, number>): number {
  if (typeof nOrprimeFactors === 'number') {
    if (nOrprimeFactors <= 0) return 0
    let res = 1
    if (!(nOrprimeFactors & 1)) {
      let e = 2
      nOrprimeFactors /= 2
      while (!(nOrprimeFactors & 1)) {
        nOrprimeFactors /= 2
        e++
      }
      res *= e
    }
    for (let f = 3; f * f <= nOrprimeFactors; f += 2) {
      if (!(nOrprimeFactors % f)) {
        let e = 2
        nOrprimeFactors /= f
        while (!(nOrprimeFactors % f)) {
          nOrprimeFactors /= f
          e++
        }
        res *= e
      }
    }
    if (nOrprimeFactors > 1) res *= 2
    return res
  }

  let res = 1
  nOrprimeFactors.forEach(count => {
    res *= count + 1
  })
  return res
}

/**
 * 返回区间 `[0, upper]` 内的所有数的约数个数.
 * @param upper 上界.
 */
function countFactorsOfAll(upper: number): Uint32Array {
  const res = new Uint32Array(upper + 1)
  for (let i = 1; i <= upper; i++) {
    for (let j = i; j <= upper; j += i) res[j]++
  }
  return res
}

/**
 * 返回约数之和.
 * @param primeFactors 质因子分解.如果分解为空，则返回 1.
 */
function sumFactors(nOrprimeFactors: number | ReadonlyMap<number, number>): number {
  if (typeof nOrprimeFactors === 'number') {
    if (nOrprimeFactors <= 0) return 0
    let res = 1
    if (!(nOrprimeFactors & 1)) {
      let cur = 1
      while (!(nOrprimeFactors & 1)) {
        nOrprimeFactors /= 2
        cur = cur * 2 + 1
      }
      res *= cur
    }
    for (let f = 3; f * f <= nOrprimeFactors; f += 2) {
      if (!(nOrprimeFactors % f)) {
        let cur = 1
        while (!(nOrprimeFactors % f)) {
          nOrprimeFactors /= f
          cur = cur * f + 1
        }
        res *= cur
      }
    }
    if (nOrprimeFactors > 1) res *= nOrprimeFactors + 1
    return res
  }

  let res = 1
  nOrprimeFactors.forEach((count, prime) => {
    let cur = 1
    for (let _ = 0; _ < count; _++) cur = cur * prime + 1
    res *= cur
  })
  return res
}

/**
 * 返回区间 `[0, upper]` 内的所有数的约数之和.
 * @param upper 上界.
 */
function sumFactorsOfAll(upper: number): number[] {
  const res = Array(upper + 1).fill(0)
  for (let i = 1; i <= upper; i++) {
    for (let j = i; j <= upper; j += i) res[j] += i
  }
  return res
}

/**
 * n 以内的最多约数个数，以及对应的最小数字.
 * @param n 上界.n <= 1e9.
 */
function maxDivisorNum(n: number): [count: number, res: number] {
  const primes = [2, 3, 5, 7, 11, 13, 17, 19, 23, 29]
  let count = 0
  let res = 1
  const dfs = (i: number, maxExp: number, curCount: number, curRes: number) => {
    if (curCount > count || (curCount === count && curRes < res)) {
      count = curCount
      res = curRes
    }
    for (let e = 1; e <= maxExp; e++) {
      curRes *= primes[i]
      if (curRes > n) break
      dfs(i + 1, e, curCount * (e + 1), curRes)
    }
  }
  dfs(0, 32 - Math.clz32(n), 1, 1)
  return [count, res]
}

/**
 * 在有 最大约数个数限制 的前提下，maxCount 最大是多少，以及对应的最小数字.
 */
function maxDivisorNumWithLimit(maxCount: number): [count: number, res: number] {
  if (maxCount === 0) return [0, 0]
  let left = 0
  let right = 1e9
  while (left <= right) {
    const mid = (left + right) >>> 1
    const [count] = maxDivisorNum(mid)
    if (count > maxCount) {
      right = mid - 1
    } else {
      left = mid + 1
    }
  }
  return maxDivisorNum(right)
}

/**
 * [min,max]以内的最多约数个数，以及对应的最小数字.
 * 1<=min<=max<=1e9
 * dfs+剪枝
 */
function maxDivisorNumInInterval(min: number, max: number): [count: number, res: number] {
  if (max - min <= 100000) {
    let count = 0
    let res = 0
    for (let i = min; i <= max; i++) {
      const curCount = countFactors(i)
      if (curCount > count) {
        count = curCount
        res = i
      }
    }
    return [count, res]
  }

  const primes = [2, 3, 5, 7, 11, 13, 17, 19, 23]
  let count = 0
  let res = 0
  const dfs = (i: number, maxExp: number, curCount: number, curRes: number) => {
    if (curRes >= min && (curCount > count || (curCount === count && curRes < res))) {
      count = curCount
      res = curRes
    }
    for (let e = 1; e <= maxExp; e++) {
      curRes *= primes[i]
      if (curRes > max) {
        break
      }
      dfs(i + 1, e, curCount * (e + 1), curRes)
    }
  }
  dfs(0, 32 - Math.clz32(max), 1, 1)
  return [count, res]
}

/**
 * n 拆分成若干连续整数的方法数/奇约数个数.
 */
function oddDivisorsNum(n: number): number {
  let res = 0
  const upper = Math.floor(Math.sqrt(n))
  for (let i = 1; i <= upper; i++) {
    if (n % i === 0) {
      if (i & 1) res++
      if (i * i < n && (n / i) & 1) res++
    }
  }
  return res
}

/**
 * 因子的中位数（偶数个因子时取小的那个）.
 */
function medianDivisor(n: number): number {
  const start = Math.floor(Math.sqrt(n))
  for (let d = start; d > 0; d--) {
    if (n % d === 0) return d
  }
  throw new Error('medianDivisor: n must be positive')
}

export {
  EratosthenesSieve,

  //
  isPrime,
  getPrimeFactors,
  countPrime,
  segmentedSieve,

  //
  getFactors,
  getFactorsOfAll,
  enumerateFactors,
  countFactors,
  countFactorsOfAll,
  sumFactors,
  sumFactorsOfAll,

  //
  maxDivisorNum,
  maxDivisorNumWithLimit,
  maxDivisorNumInInterval,
  oddDivisorsNum,
  medianDivisor
}

if (require.main === module) {
  for (let i = 1; i <= 100; i++) {
    if (countFactors(i) !== countFactors(getPrimeFactors(i))) {
      throw new Error(`countFactors(${i}) !== countFactors(getPrimeFactors(${i}))`)
    }
    if (sumFactors(i) !== sumFactors(getPrimeFactors(i))) {
      throw new Error(`sumFactors(${i}) !== sumFactors(getPrimeFactors(${i}))`)
    }
  }

  const P = new EratosthenesSieve(1e6)
  assert.strictEqual(P.isPrime(3), true)
  assert.deepStrictEqual(P.getPrimes(20), [2, 3, 5, 7, 11, 13, 17, 19])

  assert.deepStrictEqual(
    P.getPrimeFactors(20),
    new Map([
      [2, 2],
      [5, 1]
    ])
  )

  assert.deepStrictEqual(getFactors(25), [1, 5, 25])
  assert.strictEqual(countPrime(1, 1e6), P.getPrimes().length)
  assert.deepStrictEqual(
    segmentedSieve(0, 1e6),
    Array.from({ length: 1e6 }, (_, i) => i).map(v => P.isPrime(v))
  )

  // countFactors
  assert.strictEqual(
    countFactors(
      new Map([
        [2, 2],
        [3, 1]
      ])
    ),
    6
  )

  console.log(countFactors(getPrimeFactors(0)))
  console.log(sumFactors(getPrimeFactors(0)))
  console.log(countFactorsOfAll(10))
  console.log(sumFactorsOfAll(10))

  console.time('countFactors')
  getPrimeFactors(1e9)
  console.timeEnd('countFactors')

  console.log(getPrimeFactors(21))
  console.time('ffoo')
  for (let i = 0; i < 1e6; i++) {
    // getFactors(i)
    enumerateFactors(i, () => {})
  }
  console.timeEnd('ffoo')
}
