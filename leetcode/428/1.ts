export {}

function buttonWithLongestTime(events: number[][]): number {
  let max = events[0][1]
  let res = events[0][0]
  for (let i = 1; i < events.length; i++) {
    const diff = events[i][1] - events[i - 1][1]
    if (diff > max || (diff === max && events[i][0] < res)) {
      max = diff
      res = events[i][0]
    }
  }
  return res
}
