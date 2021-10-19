/**
 * @param {string[]} timePoints
 * @return {number}
 */
const findMinDifference = function (timePoints: string[]): number {
  timePoints.sort()
  const points = timePoints.map(t => {
    const [hour, min] = t.split(':')
    return Number(hour) * 60 + Number(min)
  })

  let min = Infinity
  for (let i = 1; i < points.length; i++) {
    min = Math.min(points[i] - points[i - 1], min)
    if (min == 0) return min
  }

  // 注意这里
  return Math.min(min, points[0] + 1440 - points[points.length - 1])
}

console.log(findMinDifference(['23:59', '00:00']))

export default 1
