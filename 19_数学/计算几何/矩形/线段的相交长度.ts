/**
 *
 * @param a
 * @param b
 * @param c
 * @param d
 * @returns
 * 一维线段求交集长度
 */
function getIntersectedLength(a: number, b: number, c: number, d: number) {
  // 线段相交长度
  return Math.max(0, Math.min(b, d) - Math.max(a, c))
}

if (require.main === module) {
  console.log(getIntersectedLength(1, 2, 3, 4))
}

export { getIntersectedLength }
