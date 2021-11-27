class Interval {
  constructor(public start: number, public end: number, public height: number) {}
}

/**
 * @param {number[][]} positions [最左边，边长]  1 <= positions.length <= 1000.
 * @return {number[]}
 * 返回一个堆叠高度列表 ans 。每一个堆叠高度 ans[i] 表示
 * 在通过 positions[0], positions[1], ..., positions[i] 表示的方块掉落结束后，
 * 目前所有已经落稳的方块堆叠的最高高度。
    O(n^2)解法
 */
const fallingSquares = function (positions: number[][]): number[] {
  const res: number[] = []
  const intervals: Interval[] = []

  let h = 0
  for (const pos of positions) {
    const cur = new Interval(pos[0], pos[0] + pos[1], pos[1])
    h = Math.max(h, getHeight(intervals, cur))
    res.push(h)
  }

  return res

  /**
   *
   * @param intervals
   * @param cur
   * @returns
   * 此处可用有序数据结构优化
   */
  function getHeight(intervals: Interval[], cur: Interval): number {
    let preMaxHeight = 0

    for (const interval of intervals) {
      if (interval.end <= cur.start) continue
      if (interval.start >= cur.end) continue
      preMaxHeight = Math.max(preMaxHeight, interval.height)
    }

    cur.height += preMaxHeight
    intervals.push(cur)
    return cur.height
  }
}

console.log(
  fallingSquares([
    [1, 2],
    [2, 3],
    [6, 1],
  ])
)
// 输出: [2, 5, 5]
// 解释:

// 第一个方块 positions[0] = [1, 2] 掉落：
// _aa
// _aa
// -------
// 方块最大高度为 2 。

// 第二个方块 positions[1] = [2, 3] 掉落：
// __aaa
// __aaa
// __aaa
// _aa__
// _aa__
// --------------
// 方块最大高度为5。
// 大的方块保持在较小的方块的顶部，不论它的重心在哪里，因为方块的底部边缘有非常大的粘性。

// 第三个方块 positions[1] = [6, 1] 掉落：
// __aaa
// __aaa
// __aaa
// _aa
// _aa___a
// --------------
// 方块最大高度为5。

// 因此，我们返回结果[2, 5, 5]。
export {}
