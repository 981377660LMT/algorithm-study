/* eslint-disable max-len */

/**
 * 带权(权重为等差数列)前缀和.
 * 权重首项为`start`,公差为`diff`.
 * 前缀和为: `start*a0+(start+diff)*a1+...+(start+(k-1)*diff)*ak`.
 */
function weightedPresum(arr: ArrayLike<number>, first = 1, diff = 1): (start: number, end: number) => number {
  const preSum1 = Array(arr.length + 1)
  const preSum2 = Array(arr.length + 1)
  preSum1[0] = 0
  preSum2[0] = 0
  for (let i = 0; i < arr.length; i++) {
    preSum1[i + 1] = preSum1[i] + diff * arr[i]
    preSum2[i + 1] = preSum2[i] + (first + i * diff) * arr[i]
  }

  const query = (start: number, end: number): number => {
    if (start >= end) return 0
    return preSum2[end] - preSum2[start] - start * (preSum1[end] - preSum1[start])
  }

  return query
}

export { weightedPresum }

if (require.main === module) {
  const S = weightedPresum([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
  console.log(S(0, 10))
}
