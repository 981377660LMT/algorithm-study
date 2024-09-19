/* eslint-disable no-inner-declarations */

/** 两个区间列表的交集长度. */
function intervalsIntersection(
  n1: number,
  f1: (i: number) => { left1: number; right1: number },
  n2: number,
  f2: (i: number) => { left2: number; right2: number }
): number {
  let res = 0
  let i = 0
  let j = 0
  while (i < n1 && j < n2) {
    const { left1, right1 } = f1(i)
    const { left2, right2 } = f2(j)
    if ((left1 <= right2 && right2 <= right1) || (left2 <= right1 && right1 <= right2)) {
      res += Math.min(right1, right2) - Math.max(left1, left2)
    }
    if (right1 < right2) {
      i++
    } else {
      j++
    }
  }
  return res
}

/** 遍历两个区间列表的交集. */
function enumerateIntervalsIntersection(
  n1: number,
  f1: (i: number) => { left1: number; right1: number },
  n2: number,
  f2: (i: number) => { left2: number; right2: number },
  f: (left: number, right: number, i: number, j: number) => boolean | void
): void {
  let i = 0
  let j = 0
  while (i < n1 && j < n2) {
    const { left1, right1 } = f1(i)
    const { left2, right2 } = f2(j)
    if ((left1 <= right2 && right2 <= right1) || (left2 <= right1 && right1 <= right2)) {
      if (f(Math.max(left1, left2), Math.min(right1, right2), i, j)) {
        return
      }
    }
    if (right1 < right2) {
      i++
    } else {
      j++
    }
  }
}

export { intervalsIntersection, enumerateIntervalsIntersection }

if (require.main === module) {
  // 986. 区间列表的交集
  // https://leetcode.cn/problems/interval-list-intersections/description/
  function intervalIntersection(firstList: number[][], secondList: number[][]): number[][] {
    const res: number[][] = []
    enumerateIntervalsIntersection(
      firstList.length,
      i => ({ left1: firstList[i][0], right1: firstList[i][1] }),
      secondList.length,
      i => ({ left2: secondList[i][0], right2: secondList[i][1] }),
      (left, right) => {
        res.push([left, right])
      }
    )
    return res
  }
}
