/**
 *
 * @param curRow
 * @param curCol
 * @param row
 * @param col
 * @returns
 */
function isValidPoint(curRow: number, curCol: number, row: number, col: number): boolean {
  return curRow >= 0 && curCol >= 0 && curRow < row && curCol < col
}

export { isValidPoint }
