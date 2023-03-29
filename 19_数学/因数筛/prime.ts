/* eslint-disable no-param-reassign */
/* eslint-disable no-shadow */

import assert from 'assert'

function usePrime(max: number) {
  const minPrime = new Uint32Array(max + 1)
  for (let i = 0; i <= max; i++) minPrime[i] = i

  const upper = Math.floor(Math.sqrt(max))
  for (let i = 2; i <= upper; i++) {
    if (minPrime[i] < i) continue
    for (let j = i * i; j <= max; j += i) {
      if (minPrime[j] === j) minPrime[j] = i
    }
  }

  function isPrime(n: number): boolean {
    if (n < 2) return false
    return minPrime[n] === n
  }

  /**
   * 求n的质因数分解
   *
   * @complexity log(n)
   */
  function getPrimeFactors(n: number): ReadonlyMap<number, number> {
    const res = new Map<number, number>()
    while (n > 1) {
      const p = minPrime[n]
      res.set(p, (res.get(p) || 0) + 1)
      n /= p
    }
    return res
  }

  /**
   * 求小于等于n的所有质因数
   */
  function getPrimes(n = max): readonly number[] {
    const res: number[] = []
    for (let i = 2; i <= n; i++) {
      if (i === minPrime[i]) res.push(i)
    }
    return res
  }

  return {
    isPrime,
    getPrimeFactors,
    getPrimes
  }
}

/**
 * 返回 n 的所有因子
 *
 * @complexity O(n^0.5)
 */
function getFactors(n: number): readonly number[] {
  if (n <= 0) return []

  const small: number[] = []
  const big: number[] = []

  const upper = Math.floor(Math.sqrt(n))
  for (let f = 1; f <= upper; f++) {
    if (n % f === 0) {
      small.push(f)
      big.push(n / f)
    }
  }

  if (small[small.length - 1] === big[big.length - 1]) big.pop()

  return [...small, ...big.reverse()]
}

function eulersSieve(n: number): readonly number[] {
  const primes = []
  const isPrime = new Uint8Array(n + 1).fill(1)
  for (let cand = 2; cand <= n; cand++) {
    if (isPrime[cand]) primes.push(cand)
    // prime is the smallest factor
    for (let index = 0; index < primes.length; index++) {
      const prime = primes[index]
      // overflow
      if (prime * cand > n) break
      isPrime[prime * cand] = 0
      // prime is no longer the SPF of the above multiple
      if (cand % prime === 0) break
    }
  }
  return primes
}

// O(nloglogn)
function eratosthenesSieve(n: number): readonly number[] {
  const isPrime = new Uint8Array(n + 1).fill(1)
  for (let p = 2; p * p <= n; p++) {
    if (isPrime[p]) {
      for (let j = p * p; j <= n; j += p) {
        isPrime[j] = 0
      }
    }
  }

  const primes = []
  for (let i = 2; i <= n; i++) if (isPrime[i]) primes.push(i)
  return primes
}

/**
 *
 * @param nth
 * @returns 返回第 n（从 1 开始）个质数，例如 prime(1) 返回 2。
 */
function prime(nth: number): number {
  let f = 20
  if (nth > 5e7) f = 50
  if (nth > 1e22) f = 100
  return getPrimes(f * nth)[nth - 1]
}

/**
 *
 * @param n
 * @returns 得到小于等于 n 的所有质数，返回一个数组。
 */
function getPrimes(n: number): readonly number[] {
  return n < 1000 ? eratosthenesSieve(n) : eulersSieve(n)
}

// O(n^0.5)
function isPrime(n: number): boolean {
  if (n < 2) return false
  const primes = getPrimes(Math.floor(Math.sqrt(n)))
  for (let index = 0; index < primes.length; index++) {
    const p = primes[index]
    if (n % p === 0) return false
  }
  return true
}

/**
 *
 * @param n
 * @returns 返回 n 的所有质数因子，键为质数，值为因子的指数。
 *
 */
function getPrimeFactors(n: number): ReadonlyMap<number, number> {
  const factors = new Map()
  const sqrt = Math.sqrt(n)
  for (let f = 2; f <= sqrt; f++) {
    let count = 0
    while (n % f === 0) {
      n /= f
      count++
    }
    if (count) factors.set(f, count)
  }
  if (n > 1) factors.set(n, 1)
  return factors
}

if (require.main === module) {
  const P = usePrime(1e6)
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
}

export {
  usePrime,
  eulersSieve,
  eratosthenesSieve,
  prime,
  getPrimes,
  isPrime,
  getPrimeFactors,
  getFactors
}
