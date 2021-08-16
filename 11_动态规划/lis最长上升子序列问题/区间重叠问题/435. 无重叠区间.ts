/**
 * @param {number[][]} intervals
 * @return {number}
 * 剩下的相邻区间的后一个区间的开始一定是不小于前一个区间的结束的(LIS)
 * 但是不太方便 还是贪心好
 * 按照每个区间结尾从小到大进行升序排序，优先选择结尾最短的区间，在它的后面才可能连接更多的区间
 */
const eraseOverlapIntervals = function (intervals: number[][]): number {
  if (intervals.length <= 1) return 0

  intervals.sort((a, b) => a[1] - b[1])
  let prev = intervals[0],
    remove = 0

  for (let i = 1; i < intervals.length; i++) {
    const [prevS, prevE] = prev
    const [currS, currE] = intervals[i]
    if (prevE > currS) {
      remove++
      continue
    }
    prev = intervals[i]
  }
  console.log(intervals)
  return remove
}

console.log(
  eraseOverlapIntervals([
    [1, 100],
    [11, 22],
    [1, 11],
    [2, 12],
  ])
)
export {}
