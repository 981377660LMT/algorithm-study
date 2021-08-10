// 静态极值
// 解空间就是从 0 到数组 nums 中最大最小值的差，

const solve = (arr: number[], k: number) => {
  arr.sort()
  // 小于等于这个中间值的任意两个数的差的绝对值有几个
  const countNGT = (mid: number) => {
    let l = 0
    let r = 0
    let count = 0
    while (r < arr.length) {
      while (arr[r] - arr[l] > mid) {
        l++
      }

      count += r - l
      r++
    }

    return count
  }
  // 解空间
  let l = 0
  let r = arr[arr.length - 1] - arr[0]

  while (l <= r) {
    const m = Math.floor((l + r) / 2)
    countNGT(m) >= k && (r = m - 1)
    countNGT(m) < k && (l = m + 1)
  }
  return l
}

console.log(solve([1, 5, 3, 2], 3))
