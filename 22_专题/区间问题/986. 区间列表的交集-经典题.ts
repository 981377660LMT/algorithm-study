/**
 *
 * @param firstList
 * @param secondList
 * 每个区间列表都是成对 不相交 的，并且 已经排序 。
 * 返回这 两个区间列表的交集 。
 */
function intervalIntersection(firstList: number[][], secondList: number[][]): number[][] {
  let i = 0
  let j = 0
  const res: number[][] = []
  while (i < firstList.length && j < secondList.length) {
    const [start1, end1] = firstList[i]
    const [start2, end2] = secondList[j]

    // 相交
    if (start2 <= end1 && start1 <= end2) {
      res.push([Math.max(start1, start2), Math.min(end1, end2)])
    }

    // 覆盖或相离
    if (end1 <= end2) i++
    else j++
  }

  return res
}

console.log(
  intervalIntersection(
    [
      [0, 2],
      [5, 10],
      [13, 23],
      [24, 25],
    ],
    [
      [1, 5],
      [8, 12],
      [15, 24],
      [25, 26],
    ]
  )
)
