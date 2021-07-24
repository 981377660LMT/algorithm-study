// 以k,b为键 O(n^2)
// 每次遍历产生一个map max记录值 空间复杂度就降低为O(n)
// Beware that JS’s problem when calculating big decimal numbers, use hack 1000000.0 *
const maxPoints = (points: [number, number][]) => {
  if (points.length < 2 || points == null) return points.length
  let max = 2

  for (let i = 0; i < points.length; i++) {
    let [p1x, p1y] = points[i]
    const map = new Map<string | number, number>()
    for (let j = i + 1; j < points.length; j++) {
      let [p2x, p2y] = points[j]
      let slope: number | string = (1000000.0 * (p2y - p1y)) / (p2x - p1x)
      if (!Number.isFinite(slope)) slope = 'v'
      else if (Number.isNaN(slope)) slope = 'h'
      map.set(slope, map.get(slope)! + 1 || 1)
    }
    max = Math.max(Math.max(...map.values()) + 1, max)
    console.log(map)
  }
  return max
}

console.log(
  maxPoints([
    [1, 1],
    [3, 2],
    [5, 3],
    [4, 1],
    [2, 3],
    [1, 4],
  ])
)
// 输出：4
export {}
