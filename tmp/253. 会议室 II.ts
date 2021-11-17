// 每个会议时间都会包括开始和结束的时间 intervals[i] = [starti, endi] ，
// 为避免会议冲突，同时要考虑充分利用会议室资源，请你计算至少需要多少间会议室，
// 才能满足这些会议安排。

// 1 <= intervals.length <= 104
// 0 <= starti < endi <= 106
function minMeetingRooms(intervals: number[][]): number {
  const max = Math.max(...intervals.flat())
  const bus = new Int32Array(max + 1)
  let sum = 0
  let res = 0

  for (let i = 0; i < intervals.length; i++) {
    const [start, end] = intervals[i]
    bus[start]++
    bus[end]--
  }

  for (let i = 0; i < bus.length; i++) {
    sum += bus[i]
    res = Math.max(res, sum)
  }

  return res
}

console.log(
  minMeetingRooms([
    [0, 30],
    [5, 10],
    [15, 20],
  ])
)

export {}
