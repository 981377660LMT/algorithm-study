// 给定二维空间中四点的坐标，返回四点是否可以构造一个正方形。
var validSquare = function (p1, p2, p3, p4) {
  // Distance formula: d(P, Q) = √ (x2 − x1)2 + (y2 − y1)2
  function dist(a, b) {
    return (a[0] - b[0]) * (a[0] - b[0]) + (a[1] - b[1]) * (a[1] - b[1])
  }

  // Compute the 6 pt-pt distances (squared, since we don't care about actual distance value)
  const distances = [
    dist(p1, p2),
    dist(p1, p3),
    dist(p1, p4),
    dist(p2, p3),
    dist(p2, p4),
    dist(p3, p4),
  ]
  // Sort & check for non-zero (points must be distinct), check for four equal sides, check for two equal diagonals.
  distances.sort((a, b) => a - b)

  // 边长不为0
  // 四边相等
  // 对角线相等
  return (
    distances[0] &&
    distances[0] === distances[1] &&
    distances[0] === distances[2] &&
    distances[0] === distances[3] &&
    distances[4] === distances[5]
  )
}
