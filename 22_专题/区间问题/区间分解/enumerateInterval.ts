/* eslint-disable no-inner-declarations */
/* eslint-disable max-len */

/**
 * 给定两个区间列表，每个区间列表都是成对 `不相交` 的。
 * 返回一个函数用于遍历`[allStart, allEnd)`范围内的所有区间。
 *
 * @returns
 * 返回一个函数，该函数接受一个回调函数，该回调函数接受五个参数`start`、`end`、`value1`、`value2`和`type`。
 * `start`和`end`表示当前区间的起点和终点，`value1`和`value2`表示当前区间在两个区间列表中的值，`type`表示当前区间的类型。
 * - 0: 不在两个区间列表中.
 * - 1: 在第一个区间列表中,不在第二个区间列表中.
 * - 2: 不在第一个区间列表中,在第二个区间列表中.
 * - 3: 在两个区间列表中.
 */
function enumerateInterval<V>(
  intervals1: { start: number; end: number; value: V }[],
  intervals2: { start: number; end: number; value: V }[]
): (
  allStart: number,
  allEnd: number,
  f: (
    ...args:
      | [type: '00', start: number, end: number, value1: undefined, value2: undefined]
      | [type: '10', start: number, end: number, value1: V, value2: undefined]
      | [type: '01', start: number, end: number, value1: undefined, value2: V]
      | [type: '11', start: number, end: number, value1: V, value2: V]
  ) => boolean | void
) => void {
  intervals1 = intervals1.slice().sort((a, b) => a.start - b.start)
  intervals2 = intervals2.slice().sort((a, b) => a.start - b.start)

  return (allStart, allEnd, f) => {
    let ptr1 = 0
    let ptr2 = 0
    let curStart = allStart
    while (ptr1 < intervals1.length && intervals1[ptr1].end <= curStart) ptr1++
    while (ptr2 < intervals2.length && intervals2[ptr2].end <= curStart) ptr2++

    while (curStart < allEnd) {
      const start1 = ptr1 < intervals1.length ? Math.min(intervals1[ptr1].start, allEnd) : allEnd
      const end1 = ptr1 < intervals1.length ? Math.min(intervals1[ptr1].end, allEnd) : allEnd
      const start2 = ptr2 < intervals2.length ? Math.min(intervals2[ptr2].start, allEnd) : allEnd
      const end2 = ptr2 < intervals2.length ? Math.min(intervals2[ptr2].end, allEnd) : allEnd

      // x = curStart 与两个区间相交的清况
      const intersect1 = start1 <= curStart && curStart < end1
      const intersect2 = start2 <= curStart && curStart < end2

      if (intersect1 && intersect2) {
        const minEnd = Math.min(end1, end2)
        if (f('11', curStart, minEnd, intervals1[ptr1].value, intervals2[ptr2].value)) return
        curStart = minEnd
        if (end1 === minEnd) ptr1++
        if (end2 === minEnd) ptr2++
      } else if (intersect1) {
        const curEnd = Math.min(end1, start2)
        if (f('10', curStart, curEnd, intervals1[ptr1].value, undefined)) return
        curStart = curEnd
        if (end1 === curEnd) ptr1++
      } else if (intersect2) {
        const curEnd = Math.min(end2, start1)
        if (f('01', curStart, curEnd, undefined, intervals2[ptr2].value)) return
        curStart = curEnd
        if (end2 === curEnd) ptr2++
      } else {
        const minStart = Math.min(start1, start2)
        if (f('00', curStart, minStart, undefined, undefined)) return
        curStart = minStart
      }
    }
  }
}

export { enumerateInterval }

if (require.main === module) {
  // 986. 区间列表的交集
  // https://leetcode.cn/problems/interval-list-intersections/
  function intervalIntersection(firstList: number[][], secondList: number[][]): number[][] {
    const intervals1 = firstList.map(([start, end]) => ({ start, end: end + 1, value: 0 }))
    const intervals2 = secondList.map(([start, end]) => ({ start, end: end + 1, value: 0 }))
    const res: number[][] = []
    const enumerate = enumerateInterval(intervals1, intervals2)
    enumerate(-Infinity, Infinity, (type, start, end, value1, value2) => {
      if (type === '11') {
        res.push([start, end - 1])
      }
    })
    return res
  }
}
