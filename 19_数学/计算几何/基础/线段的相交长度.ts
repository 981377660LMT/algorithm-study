/**
 *
 * @param a
 * @param b
 * @param c
 * @param d
 * @returns
 * 一维线段求交集长度
 */
const getIntersectedLength = (a: number, b: number, c: number, d: number) => {
  return Math.max(0, Math.min(b, d) - Math.max(a, c))
}

console.log(getIntersectedLength(1, 2, 3, 4))
