// 以k,b为键 O(n^2)
// Beware that JS’s problem when calculating big decimal numbers, use hack 1000000.0 *
// 注意计算的相同点
const maxPoints = (points: [number, number][]) => {
  const map = new Map<string, number>()
  const calKey = (p1: [number, number], p2: [number, number]): string => {
    let k: string
    let b: string
    if (p1[0] === p2[0]) {
      k = Infinity.toString()
      b = Infinity.toString()
    } else {
      const tmp = (p2[1] - p1[1]) / (p2[0] - p1[0])
      console.log(p1, p2, tmp)
      k = tmp.toPrecision(3)
      b = (p1[1] - tmp * p1[0]).toPrecision(3)
    }

    return k + '#' + b
  }

  for (let i = 0; i < points.length; i++) {
    for (let j = i + 1; j < points.length; j++) {
      if (i === j) continue

      const key = calKey(points[i], points[j])
      map.set(key, map.get(key)! + 1 || 1)
    }
  }
  console.log(map)
  return Math.floor(Math.sqrt([...map.values()].sort((a, b) => b - a)[0] * 2)) + 1
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
