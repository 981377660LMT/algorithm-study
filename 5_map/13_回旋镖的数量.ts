// 回旋镖 是由点 (i, j, k) 表示的元组 ，其中 i 和 j 之间的距离和 i 和 k 之间的距离相等
// i是一个枢纽 遍历每个点生成一个map 看距离有相等的没
// （需要考虑元组的顺序）
// 时间复杂度O(n^2)
// 空间复杂度O(n^2)
const numbereOfBoomerangs = (points: [number, number][]) => {
  const calDis = (p1: [number, number], p2: [number, number]) =>
    Math.sqrt((p1[0] - p2[0]) ** 2 + (p1[1] - p2[1]) ** 2)

  const map = new Map<number, Map<number, number>>()
  for (let i = 0; i < points.length; i++) {
    map.set(i, new Map())
    for (let j = 0; j < points.length; j++) {
      if (i === j) continue
      const curDis = calDis(points[i], points[j])
      if (!map.get(i)!.has(curDis)) {
        map.get(i)!.set(curDis, 1)
      } else {
        const count = map.get(i)!.get(curDis)!
        map.get(i)!.set(curDis, count + 1)
      }
    }
  }

  return [...map.values()]
    .map(map => [...map.values()])
    .flat()
    .filter(num => num >= 2)
    .reduce((pre, cur) => pre + cur * (cur - 1), 0)
}

console.log(
  numbereOfBoomerangs([
    [0, 0],
    [1, 0],
    [-1, 0],
    [0, 1],
    [0, -1],
  ])
)

export {}
