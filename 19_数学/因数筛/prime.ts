function eulersSieve(n: number): number[] {
  const primes = []
  const isPrime = Array(n + 1).fill(true)
  for (let cand = 2; cand <= n; cand++) {
    if (isPrime[cand]) primes.push(cand)
    // prime is the smallest factor
    for (const prime of primes) {
      // overflow
      if (prime * cand > n) break
      isPrime[prime * cand] = false
      // prime is no longer the SPF of the above multiple
      if (cand % prime === 0) break
    }
  }
  return primes
}

// O(nloglogn)
function eratosthenesSieve(n: number): number[] {
  const isPrime = Array(n + 1).fill(true)
  for (let p = 2; p * p <= n; p++) {
    if (isPrime[p]) {
      for (let j = p * p; j <= n; j += p) {
        isPrime[j] = false
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
  return primesLeq(f * nth)[nth - 1]
}

/**
 *
 * @param n
 * @returns 得到小于等于 n 的所有质数，返回一个数组。
 */
function primesLeq(n: number): number[] {
  return n < 1000 ? eratosthenesSieve(n) : eulersSieve(n)
}

// O(n^0.5)
function isPrime(n: number): boolean {
  if (n < 2) return false
  const primes = primesLeq(Math.floor(Math.sqrt(n)))
  for (const p of primes) if (n % p === 0) return false
  return true
}

/**
 *
 * @param n
 * @returns 返回 n 的所有质数因子，键为质数，值为因子的指数。
 */
function primeFactors(n: number): Map<number, number> {
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

/**
 * @param n
 * @returns 返回 n 的所有因子
 */
function factors(n: number): number[] {
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

if (require.main === module) {
  console.log(primeFactors(20))
  console.log(factors(25))
  console.log(factors(1))
}

export { eulersSieve, eratosthenesSieve, prime, primesLeq, isPrime, primeFactors, factors }
