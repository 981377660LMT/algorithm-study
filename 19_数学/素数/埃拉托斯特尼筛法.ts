// 埃拉托斯特尼筛法（sieve of eratosthenes）生成素数序列
// 其原理是剔除所有可被素数整除的非素数
/**
 *
 * @param n  统计所有小于非负整数 n 的质数的数量。
 */
function countPrimes(n: number): number {
  const visited = Array<boolean>(n).fill(false)
  let res = 0

  for (let i = 2; i < n; i++) {
    if (visited[i]) continue
    res++
    for (let j = i * i; j < n; j += i) {
      visited[j] = true
    }
  }

  return res
}

console.log(countPrimes(10))
// 运用比特表（Bitmap）算法对筛法进行内存优化
