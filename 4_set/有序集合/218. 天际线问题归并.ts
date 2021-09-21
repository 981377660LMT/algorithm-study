/**
 * @param {number[][]} buildings  左,右,高
 * @return {number[][]}
 * 求轮廓中的所有的水平线的左端点
 *  * 两两合并，最后一起合并
 * https://leetcode-cn.com/problems/the-skyline-problem/solution/fen-er-zhi-zhi-er-fen-fa-dui-by-powcai/

 */
var getSkyline = function (buildings: number[][]): number[][] {
  const merge = (left: number[][], right: number[][]): number[][] => {
    console.log(1, left, right)
    const res: number[][] = []
    // 当前左数组里建筑物高度
    let lheight = 0
    // 当前右数组里建筑物高度
    let rheight = 0
    let i = 0
    let j = 0
    let curPoint: number[]

    while (i < left.length && j < right.length) {
      if (left[i][0] < right[j][0]) {
        curPoint = [left[i][0], Math.max(left[i][1], rheight)]
        lheight = left[i][1]
        i++
      } else if (left[i][0] > right[j][0]) {
        curPoint = [right[j][0], Math.max(right[j][1], lheight)]
        rheight = right[j][1]
        j++
      } else {
        curPoint = [left[i][0], Math.max(left[i][1], right[j][1])]
        lheight = left[i][1]
        rheight = right[j][1]
        i++
        j++
      }
      if (res.length === 0 || res[res.length - 1][1] !== curPoint[1]) res.push(curPoint)
    }

    res.push(...left.slice(i), ...right.slice(j))
    return res
  }

  if (!buildings.length) return []
  if (buildings.length === 1)
    return [
      [buildings[0][0], buildings[0][2]],
      [buildings[0][1], 0],
    ]

  const mid = buildings.length >> 1
  const left = buildings.slice(0, mid)
  const right = buildings.slice(mid)
  return merge(getSkyline(left), getSkyline(right))
}

getSkyline([
  [2, 9, 10],
  [3, 7, 15],
  [5, 12, 12],
  [15, 20, 10],
  [19, 24, 8],
])

export {}
