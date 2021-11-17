// 请你判断一个人是否能够参加这里面的全部会议。
function canAttendMeetings(intervals: number[][]): boolean {
  if (intervals.length <= 1) return true

  intervals.sort((a, b) => a[0] - b[0])

  for (let i = 1; i < intervals.length; i++) {
    const [_preStart, preEnd] = intervals[i - 1]
    const [curStart, _curEnd] = intervals[i]
    if (curStart < preEnd) return false
  }

  return true
}

console.log(
  canAttendMeetings([
    [0, 30],
    [5, 10],
    [15, 20],
  ])
)
