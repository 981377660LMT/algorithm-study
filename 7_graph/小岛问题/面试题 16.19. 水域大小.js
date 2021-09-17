/**
 * @param {number[][]} land
 * @return {number[]}
 */
const pondSizes = function (land) {
  let res = []
  let area = 0
  // 逆时针8个方向
  let dirs = [
    [-1, 0],
    [-1, 1],
    [0, 1],
    [1, 1],
    [1, 0],
    [1, -1],
    [0, -1],
    [-1, -1],
  ]

  const dfs = (i, j) => {
    land[i][j] = -1
    area += 1
    for (let [dx, dy] of dirs) {
      let x = i + dx,
        y = j + dy
      if (x >= 0 && x < land.length && y >= 0 && y < land[0].length && land[x][y] === 0) {
        dfs(x, y)
      }
    }
  }

  for (let i = 0; i < land.length; i++) {
    for (let j = 0; j < land[i].length; j++) {
      if (land[i][j] === 0) {
        dfs(i, j)
        res.push(area)
        area = 0
      }
    }
  }

  return res
}
