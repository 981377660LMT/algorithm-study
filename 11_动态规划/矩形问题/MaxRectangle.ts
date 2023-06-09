/* eslint-disable no-inner-declarations */

// 直方图中的最大矩形

/**
 * 直方图中的最大矩形.
 * @param heights 直方图的高度.
 * @param f (start, end, height) : 以 [start, end) 为底, 高度为 height 的矩形.
 */
function maxRectangleInHistogram(
  heights: number[],
  f: (start: number, end: number, height: number) => void
): void {
  const n = heights.length
  const stack: [number, number][] = []
  for (let right = 0; right <= n; right++) {
    const rightHeight = right < n ? heights[right] : 0
    let j = right
    while (stack.length) {
      const [left, leftHeight] = stack[stack.length - 1]
      if (leftHeight < rightHeight) break
      f(left, right, leftHeight)
      stack.pop()
      j = left
    }
    stack.push([j, rightHeight])
  }
}

/**
 * 矩阵中的最大矩形.
 * @param grid 二维矩阵."1"或者1表示有效区域,"0"或者0表示无效区域.
 * @param f (r1, r2, c1, c2) : `[r1,r2) x [c1,c2)`区域.
 */
function maxRectangle1(
  grid: ArrayLike<ArrayLike<string | number | boolean>>,
  f: (r1: number, r2: number, c1: number, c2: number) => void
): void {
  const ROW = grid.length
  if (!ROW) return

  const COL = grid[0].length
  const heights = new Uint32Array(COL)
  const zero = new Uint32Array(COL + 1)
  for (let i = 0; i < ROW; i++) {
    const row = grid[i]
    const cache = i + 1 !== ROW ? grid[i + 1] : []
    for (let c = 0; c < COL; c++) {
      heights[c] = (heights[c] + 1) * +row[c]
      if (i + 1 !== ROW) zero[c + 1] = zero[c] + (+cache[c] ^ 1)
      else zero[c + 1] = zero[c] + 1
    }

    // maxRectangleHistogram
    const stack: [number, number][] = []
    for (let right = 0; right <= COL; right++) {
      const rightHeight = right < COL ? heights[right] : 0
      let pos = right
      while (stack.length) {
        const [left, leftHeight] = stack[stack.length - 1]
        if (leftHeight < rightHeight) break
        if (leftHeight && zero[right] - zero[left]) f(i - leftHeight + 1, i + 1, left, right)
        stack.pop()
        pos = left
      }
      stack.push([pos, rightHeight])
    }
  }
}

/**
 * 矩阵中的最大矩形.
 * @param grid 二维矩阵."1"或者1表示有效区域,"0"或者0表示无效区域.
 * @param f (r1, r2, c1, c2) : `[r1,r2) x [c1,c2)`区域.
 * @returns (maxArea, maxRect): maxArea: 最大矩形的面积; maxRect: `[r1,r2) x [c1,c2)`区域.
 */
function maxRectangle2(
  grid: ArrayLike<ArrayLike<string | number | boolean>>
): [maxArea: number, maxRectangle: [r1: number, r2: number, c1: number, c2: number]] {
  let maxArea = 0
  let maxRect: [number, number, number, number] = [0, 0, 0, 0]

  function f(r1: number, r2: number, c1: number, c2: number): void {
    const area = (r2 - r1) * (c2 - c1)
    if (area > maxArea) {
      maxArea = area
      maxRect = [r1, r2, c1, c2]
    }
  }

  maxRectangle1(grid, f)
  return [maxArea, maxRect]
}

export { maxRectangle1, maxRectangle2, maxRectangleInHistogram }
if (require.main === module) {
  // https://leetcode.cn/problems/largest-rectangle-in-histogram/
  function largestRectangleArea(heights: number[]): number {
    let res = 0
    maxRectangleInHistogram(heights, (start, end, height) => {
      res = Math.max(res, height * (end - start))
    })
    return res
  }

  // https://leetcode.cn/problems/maximal-rectangle/
  function maximalRectangle(matrix: string[][]): number {
    return maxRectangle2(matrix)[0]
  }
}
