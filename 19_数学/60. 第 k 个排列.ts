/**
 * @param {number} n
 * @param {number} k  1 <= k <= n!
 * @return {string}
 * 给出集合 [1,2,3,...,n]，其所有元素共有 n! 种排列
 */
const getPermutation = function (n: number, k: number): string {
  const factorial = (n: number): number => (n <= 1 ? 1 : factorial(n - 1) * n)
  const numbers = Array.from<number, number>({ length: n }, (_, i) => i + 1)
  const res: number[] = []
  k-- // 序号从0计算
  while (n) {
    n--
    const index = ~~(k / factorial(n))
    k = k % factorial(n)
    res.push(numbers[index])
    numbers.splice(index, 1)
  }

  return res.join('')
}

console.log(getPermutation(4, 9))

export {}
