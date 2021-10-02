// 1 <= ranges.length <= 50
// 1 <= starti <= endi <= 50

import { BIT } from './BIT'

// 1 <= left <= right <= 50
function isCovered(ranges: number[][], left: number, right: number): boolean {
  const bit = new BIT(51)
  const visited = new Set<number>() // 让每个点只添加一次

  for (const [l, r] of ranges) {
    for (let i = l; i <= r; i++) {
      if (visited.has(i)) continue
      visited.add(i)
      bit.add(i, 1)
    }
  }

  return bit.sumRange(left, right) === right - left + 1
}

console.log(
  isCovered(
    [
      [1, 2],
      [3, 4],
      [5, 6],
    ],
    2,
    5
  )
)
// 输出：true
// 解释：2 到 5 的每个整数都被覆盖了：
// - 2 被第一个区间覆盖。
// - 3 和 4 被第二个区间覆盖。
// - 5 被第三个区间覆盖。
