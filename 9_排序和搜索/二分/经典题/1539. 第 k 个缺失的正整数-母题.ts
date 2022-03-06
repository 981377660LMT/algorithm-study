// 给你一个 严格升序排列 的正整数数组 arr 和一个整数 k 。
// 请你找到这个数组里第 k 个缺失的正整数。
function findKthMex(arr: number[], k: number): number {
  let l = 0
  let r = arr.length - 1
  while (l <= r) {
    const mid = (l + r) >> 1
    const missing = arr[mid] - (mid + 1)
    if (missing >= k) r = mid - 1
    else l = mid + 1
  }

  return k + l
}

if (require.main === module) {
  console.log(
    findKthMex(
      [
        96, 44, 99, 25, 61, 84, 88, 18, 19, 33, 60, 86, 52, 19, 32, 47, 35, 50, 94, 17, 29, 98, 22,
        21, 72, 100, 40, 84,
      ].sort((a, b) => a - b),
      35
    )
  )
}

export { findKthMex }
// 解释：缺失的正整数包括 [1,5,6,8,9,10,12,13,...] 。第 5 个缺失的正整数为 9 。
