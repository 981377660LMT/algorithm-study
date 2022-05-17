type Interval = [number, number]

/**
 * @param {Interval[][]} schedules
 * @return {Interval[]}
 * [start, end]含有0～24的整数，意味着当前时间段已经有安排。
 */
function findMeetingSlots(schedules: Interval[][]): Interval[] {
  const times = schedules.flat().sort((a, b) => a[0] - b[0])
  // [ [ 8, 9 ], [ 10, 13 ], [ 11, 12 ], [ 13, 15 ], [ 13, 18 ] ]
  const res: Interval[] = []
  let preEnd = 0
  times.forEach(([start, end]) => {
    if (preEnd < start) res.push([preEnd, start])
    preEnd = Math.max(end, preEnd)
  })

  if (preEnd !== 24) res.push([preEnd, 24])

  return res
}

console.log(
  findMeetingSlots([
    [
      [13, 15],
      [11, 12],
      [10, 13],
    ], //成员1的安排
    [[8, 9]], //成员2的安排
    [[13, 18]], //成员3的安排
  ])
)

// [[0,8],[9,10],[18,24]]

export {}
