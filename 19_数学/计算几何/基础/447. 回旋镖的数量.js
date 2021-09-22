/**
 * @param {number[][]} points
 * @return {number}   返回平面上所有回旋镖的数量
 */
var numberOfBoomerangs = function (points) {
  const calDis = (p1, p2) => Math.sqrt((p1[0] - p2[0]) ** 2 + (p1[1] - p2[1]) ** 2)

  const map = new Map()
  for (let i = 0; i < points.length; i++) {
    map.set(i, new Map())
    for (let j = 0; j < points.length; j++) {
      if (i === j) continue
      const curDis = calDis(points[i], points[j])
      if (!map.get(i).has(curDis)) {
        map.get(i).set(curDis, 1)
      } else {
        const count = map.get(i).get(curDis)
        map.get(i).set(curDis, count + 1)
      }
    }
  }

  return [...map.values()]
    .flatMap(map => [...map.values()])
    .reduce((pre, cur) => pre + cur * (cur - 1), 0)
}
