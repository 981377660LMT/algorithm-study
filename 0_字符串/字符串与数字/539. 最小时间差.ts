/**
 * @param {string[]} timePoints
 * @return {number}
 */
var findMinDifference = function (timePoints: string[]): number {
  timePoints.sort()
  let minPoints = timePoints.map(t => Number(t[0] + t[1]) * 60 + Number(t[3] + t[4]))

  let min = Infinity
  for (let i = 1; i < minPoints.length; i++) {
    min = Math.min(minPoints[i] - minPoints[i - 1], min)
    if (min == 0) return min
  }

  // 注意这里
  return Math.min(min, minPoints[0] + 1440 - minPoints[minPoints.length - 1])
}

console.log(findMinDifference(['23:59', '00:00']))

export default 1
