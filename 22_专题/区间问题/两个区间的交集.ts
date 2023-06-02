/**
 * 获取两个区间的交集.
 */
function getIntersect(
  interval1: [number, number],
  interval2: [number, number]
): [number, number] | null {
  if (interval1[0] > interval2[1] || interval2[0] > interval1[1]) {
    return null
  }
  return [Math.max(interval1[0], interval2[0]), Math.min(interval1[1], interval2[1])]
}
