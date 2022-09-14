/* eslint-disable @typescript-eslint/no-unused-vars */
/**
 * @param {number} ax1
 * @param {number} ay1
 * @param {number} ax2
 * @param {number} ay2
 * @param {number} bx1
 * @param {number} by1
 * @param {number} bx2
 * @param {number} by2
 * @return {number}
 * 给你 二维 平面上两个 由直线构成的 矩形，
 * !请你计算并返回两个矩形覆盖的总面积。
 * 每个矩形由其 左下 顶点和 右上 顶点坐标表示：
 * 求相交矩形的面积
 * 矩形相交
 */
function computeArea(ax1, ay1, ax2, ay2, bx1, by1, bx2, by2) {
  // x轴投影相交线段长度*y轴投影相交线段长度 就是 重叠面积
  const x = calIntersectedLength(ax1, ax2, bx1, bx2)
  const y = calIntersectedLength(ay1, ay2, by1, by2)
  const overlap = x * y // 重叠面积
  return (ay2 - ay1) * (ax2 - ax1) + (by2 - by1) * (bx2 - bx1) - overlap

  // 线段相交长度
  function calIntersectedLength(a, b, c, d) {
    return Math.max(0, Math.min(b, d) - Math.max(a, c))
  }
}
