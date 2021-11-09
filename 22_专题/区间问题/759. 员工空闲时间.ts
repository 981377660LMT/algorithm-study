class Interval {
  start: number
  end: number
  constructor(start: number, end: number) {
    this.start = start
    this.end = end
  }
}
/**
 * 
 * @param schedule 
 * @description
 * 给定员工的 schedule 列表，表示每个员工的工作时间。
   每个员工都有一个非重叠的时间段  Intervals 列表，这些时间段已经排好序。
   返回表示 所有 员工的 共同，正数长度的空闲时间 的有限时间段的列表，同样需要排好序。
 */
function employeeFreeTime(schedule: Interval[][]): Interval[] {}

export {}
// 输入：schedule = [[[1,2],[5,6]],[[1,3]],[[4,10]]]
// 输出：[[3,4]]
// 解释：
// 共有 3 个员工，并且所有共同的
// 空间时间段是 [-inf, 1], [3, 4], [10, inf]。
// 我们去除所有包含 inf 的时间段，因为它们不是有限的时间段。
