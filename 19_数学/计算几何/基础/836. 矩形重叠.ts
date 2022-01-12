import { getIntersectedLength } from './线段的相交长度'

/**
 * @param {number[]} rec1
 * @param {number[]} rec2
 * @return {boolean}
 */
function isRectangleOverlap(rec1: number[], rec2: number[]): boolean {
  const [x1, y1, x2, y2] = rec1
  const [x3, y3, x4, y4] = rec2
  return !!getIntersectedLength(x1, x2, x3, x4) && !!getIntersectedLength(y1, y2, y3, y4)
}

console.log(isRectangleOverlap([0, 0, 2, 2], [1, 1, 3, 3]))
