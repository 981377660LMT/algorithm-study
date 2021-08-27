// 你需要设计一个平滑器来让每一个单元的灰度成为平均灰度 (向下舍入) ，
// 平均灰度的计算是周围的8个单元和它本身的值求平均，如果周围的单元格不足八个，
// 则尽可能多的利用它们。
/**
 * @param {number[][]} img
 * @return {number[][]}
 */
var imageSmoother = img => {
  let rows = img.length,
    cols = img[0].length
  let ret = new Array(rows).fill(0).map(_ => new Array(cols).fill(0))
  for (let r = 0; r < rows; ++r) {
    for (let c = 0; c < cols; ++c) {
      let count = 0
      for (let x of [-1, 0, 1])
        for (let y of [-1, 0, 1])
          if (isValid(r + x, c + y, rows, cols)) {
            count++
            ret[r][c] += img[r + x][c + y]
          }
      ret[r][c] = Math.floor(ret[r][c] / count)
    }
  }
  return ret
}

const isValid = (r, c, rows, cols) => r < rows && r >= 0 && c < cols && c >= 0
