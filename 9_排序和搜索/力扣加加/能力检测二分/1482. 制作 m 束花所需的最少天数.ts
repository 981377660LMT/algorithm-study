/**
 *
 * @param bloomDay  1 <= n <= 10^5
 * @param m  1 <= m <= 10^6
 * @param k
 * @description
 * 现需要制作 m 束花。制作花束时，需要使用花园中 `相邻的 k 朵花` 。
 * 花园中有 n 朵花，第 i 朵花会在 bloomDay[i] 时盛开，恰好 可以用于 一束 花中。
 * 请你返回从花园中摘 m 束花需要等待的最少的天数。如果不能摘到 m 束花则返回 -1 。
 */
function minDays(bloomDay: number[], m: number, k: number): number {
  if (m * k > bloomDay.length) return -1

  let left = 1
  let right = Math.max(...bloomDay)
  while (left <= right) {
    const waitDay = (left + right) >> 1
    if (canSatisfy(waitDay)) right = waitDay - 1
    else left = waitDay + 1
  }

  return left

  function canSatisfy(waitDay: number): boolean {
    let flowerCount = 0
    let bouquetCount = 0

    for (const bloom of bloomDay) {
      if (bloom > waitDay) {
        flowerCount = 0
      } else {
        flowerCount++
        if (flowerCount >= k) {
          bouquetCount++
          flowerCount = 0
        }
      }

      // 加速
      if (bouquetCount >= m) return true
    }

    return bouquetCount >= m
  }
}

console.log(minDays([1, 10, 3, 10, 2], 3, 1))
