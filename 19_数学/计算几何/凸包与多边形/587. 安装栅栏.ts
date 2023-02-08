/**
 * @param {number[][]} trees
 * @return {number[][]}
 * 用最短的绳子围起所有的树
   找到正好位于栅栏边界上的树的坐标
   @summary
   寻找凸包:单调链
   凸多边形是指所有内角大小都在[0, \pi][0,π]范围内的简单多边形。
   在平面上能包含所有给定点的最小凸多边形叫做凸包。
   @description
    对于每个新的点，我们检查当前点是否在最后两个点的逆时针方向上。
    如果是的话，当前点直接被压入凸壳 hull 中，
    如果不是的话（即 orientation 返回的结果为正数），
    我们可以知道栈顶的元素在凸壳里面而不是凸壳边上。
    @link
    https://leetcode-cn.com/problems/erect-the-fence/solution/an-zhuang-zha-lan-by-leetcode/
    时间复杂度：O(nlog(n))。
  */
const outerTrees = (points: [x: number, y: number][]) => {
  if (points.length <= 1) return points
  points.sort((a, b) => a[0] - b[0] || a[1] - b[1])
  const cand: [x: number, y: number][] = []

  // 寻找凸壳的下半部分
  for (let i = 0; i < points.length; i++) {
    while (
      cand.length >= 2 &&
      // 顺时针出栈
      cross(cand[cand.length - 2], cand[cand.length - 1], points[i]) < 0
    ) {
      cand.pop()
    }
    cand.push(points[i])
  }

  // 目前求解出来的部分只包括凸壳的下半部分。现在我们需要求出凸壳的上半部分
  // result.pop()
  for (let i = points.length - 1; i >= 0; i--) {
    while (
      cand.length >= 2 &&
      // 顺时针出栈
      cross(cand[cand.length - 2], cand[cand.length - 1], points[i]) < 0
    ) {
      cand.pop()
    }
    cand.push(points[i])
  }

  return [...new Set(cand)]

  /**
   *
   * @param a
   * @param b
   * @param c
   * @returns
   * 大于0 则 abc逆时针
   * 小于0 则 abc顺时针
   * ab 叉乘 ac
   */
  function cross(a: number[], b: number[], c: number[]) {
    const ab = [b[0] - a[0], b[1] - a[1]]
    const ac = [c[0] - a[0], c[1] - a[1]]
    return ab[0] * ac[1] - ab[1] * ac[0]
  }
}

console.log(
  outerTrees([
    [1, 1],
    [2, 2],
    [2, 0],
    [2, 4],
    [3, 3],
    [4, 2]
  ])
)
