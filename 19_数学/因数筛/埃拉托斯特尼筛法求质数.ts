// 埃拉托斯特尼筛法（sieve of eratosthenes）生成素数序列
// 其原理是剔除所有可被素数整除的非素数
// # 筛法的复杂度为n/1 + n/2 + n/3 + … + n/n 渐进为O(n * logn)
// # 而gcd的复杂度为O(logn)，所以总复杂度为O(n * logn * logn)。

/**
 *
 * @param n  统计所有不超过非负整数 n 的质数的数量。
 */
function countPrimes(n: number): number {
  const isPrime = new Uint8Array(n + 1)
  let res = 0

  for (let f = 2; f <= n; f++) {
    if (isPrime[f]) continue
    res++
    for (let j = f * f; j <= n; j += f) {
      isPrime[j] = 1
    }
  }

  return res
}

// console.log(countPrimes(10))
// 运用比特表（Bitmap）算法对筛法进行内存优化

/**
 *
 * @param upper
 * @returns  不超过upper的质数
 */
function getPrimes(upper: number): number[] {
  const visited = Array<boolean>(upper + 1).fill(false)

  for (let i = 2; i <= upper; i++) {
    if (visited[i]) continue
    for (let j = i * i; j <= upper; j += i) {
      visited[j] = true
    }
  }

  const res: number[] = []
  for (let i = 2; i <= upper; i++) {
    if (!visited[i]) res.push(i)
  }

  return res
}

export { getPrimes, countPrimes }
