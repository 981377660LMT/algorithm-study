// 如果它们有交点，请计算其交点，没有交点则返回空值。
// 若有多个交点（线段重叠）则返回 X 值最小的点，X 坐标相同则返回 Y 值最小的点。
// https://leetcode-cn.com/problems/intersection-lcci/solution/wo-jue-de-wo-yi-jing-hen-nu-li-liao-ke-yi-jiao-zhu/

/**
 *
 * @param start1
 * @param end1
 * @param start2
 * @param end2
 * @summary
 * 将线段表示出 crammer法则求交点即可
 * Cramer 法则只在最后一步使用除法，中途不会出现除零和精度的错误
 */
function intersection(
  start1: number[],
  end1: number[],
  start2: number[],
  end2: number[]
): number[] {
  const det = (a: number, b: number, c: number, d: number) => a * d - b * c
  const [x1, y1] = start1
  const [x2, y2] = end1
  const [x3, y3] = start2
  const [x4, y4] = end2
  const Up1 = det(x4 - x2, x4 - x3, y4 - y2, y4 - y3)
  const Up2 = det(x1 - x2, x4 - x2, y1 - y2, y4 - y2)
  const Down = det(x1 - x2, x4 - x3, y1 - y2, y4 - y3)

  // 唯一解
  if (Down) {
    const lam = Up1 / Down
    const eta = Up2 / Down
    const isValid = lam >= 0 && lam <= 1 && eta >= 0 && eta <= 1
    if (isValid) return [lam * x1 + (1 - lam) * x2, lam * y1 + (1 - lam) * y2]
    else return []
  }

  // 无解 对应直线平行
  if (Up1 || Up2) return []

  // 无穷解 对应直线重合
  // groupi[0]的点更加左
  const group1 = [start1, end1].sort((a, b) => a[0] - b[0] || a[1] - b[1])
  const group2 = [start2, end2].sort((a, b) => a[0] - b[0] || a[1] - b[1])
  // 这是什么比较？表示线段不重合(数组会化成数字比较)
  if (group1[1] < group2[0] || group2[1] < group1[0]) return []

  return group1[0] < group2[0] ? group2[0] : group1[0]
}

console.log(intersection([0, 0], [1, 0], [1, 1], [0, -1]))
