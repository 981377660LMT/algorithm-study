function removeInterval(intervals: number[][], toBeRemoved: number[]): number[][] {
  const [left, right] = toBeRemoved
  const res: [number, number][] = []

  for (const [start, end] of intervals) {
    // 不重叠
    if (end < left || start > right) {
      res.push([start, end])
    } else {
      // 重叠
      if (start < left) res.push([start, left])
      if (right < end) res.push([right, end])
    }
  }

  return res
}

console.log(
  removeInterval(
    [
      [0, 2],
      [3, 4],
      [5, 7],
    ],
    [1, 6]
  )
)
// 输出：[[0,1],[6,7]]

console.log(
  removeInterval(
    [
      [-5, -4],
      [-3, -2],
      [1, 2],
      [3, 5],
      [8, 9],
    ],
    [-1, 4]
  )
)
