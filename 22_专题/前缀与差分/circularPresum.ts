/**
 * 环形数组前缀和.
 * @param arr 环形数组的循环部分.
 * @returns 返回区间 `[start, end)` 的和.
 */
function cicularPreSum(arr: ArrayLike<number>): (start: number, end: number) => number {
  const n = arr.length
  const preSum = Array(n + 1).fill(0)
  for (let i = 0; i < n; i++) preSum[i + 1] = preSum[i] + arr[i]
  const cal = (r: number) => preSum[n] * Math.floor(r / n) + preSum[r % n]
  const query = (start: number, end: number) => {
    if (start >= end) return 0
    return cal(end) - cal(start)
  }
  return query
}

export { cicularPreSum }

if (require.main === module) {
  const query = cicularPreSum([1, 2, 3, 4, 5])
  console.log(query(20, 11))
}
