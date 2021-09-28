// 埃拉托斯特尼筛法（sieve of eratosthenes）生成素数序列
// 其原理是剔除所有可被素数整除的非素数
/**
 *
 * @param n  统计所有不超过非负整数 n 的质数的数量。
 */
function countPrimes(n: number): number {
  const visited = Array<boolean>(n + 1).fill(false)
  let res = 0

  for (let i = 2; i <= n; i++) {
    if (visited[i]) continue
    res++
    for (let j = i * i; j <= n; j += i) {
      visited[j] = true
    }
  }

  return res
}

// console.log(countPrimes(10))
// 运用比特表（Bitmap）算法对筛法进行内存优化

/**
 *
 * @param upper  找出不超过upper的质数
 * @returns
 */
function getPrimes(upper: number): number[] {
  const visited = Array<boolean>(upper + 1).fill(false)
  const res: number[] = []

  for (let i = 2; i <= upper; i++) {
    if (visited[i]) continue
    res.push(i)
    for (let j = i * i; j <= upper; j += i) {
      visited[j] = true
    }
  }

  return res
}

export { getPrimes, countPrimes }
