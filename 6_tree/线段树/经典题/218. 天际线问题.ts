function getSkyline(buildings: number[][]): number[][] {}

console.log(
  getSkyline([
    [2, 9, 10],
    [3, 7, 15],
    [5, 12, 12],
    [15, 20, 10],
    [19, 24, 8],
  ])
)
// 输出：[[2,10],[3,15],[7,12],[12,0],[15,10],[20,8],[24,0]]
