export {}

/* eslint-disable prefer-destructuring */
interface Rectangle {
  sx: number
  sy: number
  ex: number
  ey: number
}

function checkValidCuts(n: number, rectangles: number[][]): boolean {
  const arr: Rectangle[] = rectangles.map(r => ({
    sx: r[0],
    sy: r[1],
    ex: r[2],
    ey: r[3]
  }))

  const sorted1 = [...arr].sort((a, b) => {
    if (a.sx !== b.sx) return a.sx - b.sx
    return a.sy - b.sy
  })
  if (check(sorted1, true)) return true
  const sorted2 = [...arr].sort((a, b) => {
    if (a.sy !== b.sy) return a.sy - b.sy
    return a.sx - b.sx
  })
  if (check(sorted2, false)) return true
  return false

  function check(sortedRects: Rectangle[], isVertical: boolean): boolean {
    const length = sortedRects.length
    if (length < 3) return false

    const preMaxEnd: number[] = Array(length).fill(0)
    const sufMinStart: number[] = Array(length).fill(0)

    preMaxEnd[0] = isVertical ? sortedRects[0].ex : sortedRects[0].ey
    for (let i = 1; i < length; i++) {
      const currentEnd = isVertical ? sortedRects[i].ex : sortedRects[i].ey
      preMaxEnd[i] = Math.max(preMaxEnd[i - 1], currentEnd)
    }
    sufMinStart[length - 1] = isVertical ? sortedRects[length - 1].sx : sortedRects[length - 1].sy
    for (let i = length - 2; i >= 0; i--) {
      const currentStart = isVertical ? sortedRects[i].sx : sortedRects[i].sy
      sufMinStart[i] = Math.min(sufMinStart[i + 1], currentStart)
    }

    const splits: number[] = []
    for (let i = 0; i < length - 1; i++) {
      if (preMaxEnd[i] <= sufMinStart[i + 1]) {
        splits.push(i)
      }
    }
    return splits.length >= 2
  }
}
