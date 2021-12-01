/**
 * @param {number[][]} mat
 * @return {number[]}
 * @summary The key here is to realize that the sum of indices on all diagonals are equal.
 * 1 <= nums.length <= 10^5
   1 <= nums[i].length <= 10^5
 */
const findDiagonalOrder = function (mat: number[][]): number[] {
  const record = new Map<number, number[]>()
  for (let i = 0; i < mat.length; i++) {
    for (let j = 0; j < mat[i].length; j++) {
      const key = i + j
      if (!record.has(key)) record.set(key, [])
      record.get(key)!.push(mat[i][j])
    }
  }

  return [...record.values()].flatMap(row => row.reverse())
}

console.log(findDiagonalOrder([[1, 2, 3, 4, 5], [6, 7], [8], [9, 10, 11], [12, 13, 14, 15, 16]]))
