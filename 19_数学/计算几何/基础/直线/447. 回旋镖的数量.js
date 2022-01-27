/**
 * @param {number[][]} points
 * @return {number}   返回平面上所有回旋镖的数量
 * 我们可以枚举每个 points[i]，将其当作 V 型的拐点。
 */
function numberOfBoomerangs(points) {
  let res = 0
  for (const p1 of points) {
    const counter = new Map()
    for (const p2 of points) {
      const dis = (p1[0] - p2[0]) * (p1[0] - p2[0]) + (p1[1] - p2[1]) * (p1[1] - p2[1])
      counter.set(dis, (counter.get(dis) || 0) + 1)
    }
    for (const [_, m] of counter.entries()) {
      res += m * (m - 1)
    }
  }
  return res
}
