// 左闭右开区间
function removeInterval(intervals: number[][], toBeRemoved: number[]): number[][] {
  const [badStart, badEnd] = toBeRemoved
  const res: [number, number][] = []

  intervals.forEach(([start, end]) => {
    // 不重叠
    if (end < badStart || start > badEnd) {
      res.push([start, end])
    } else {
      if (start < badStart) res.push([start, badStart]) // 重叠
      if (badEnd < end) res.push([badEnd, end])
    }
  })

  return res
}

console.log(
  removeInterval(
    [
      [0, 2],
      [3, 4],
      [5, 7]
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
      [8, 9]
    ],
    [-1, 4]
  )
)
