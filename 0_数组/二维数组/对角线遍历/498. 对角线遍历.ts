/**
 * @param {number[][]} mat
 * @return {number[]}
 * @summary The key here is to realize that the sum of indices on all diagonals are equal.
 * 1 <= m, n <= 104
   1 <= m * n <= 104
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

  const res: number[] = []
  for (const [key, nums] of record.entries()) {
    key & 1 ? res.push(...nums) : res.push(...nums.reverse())
  }

  console.table(record)
  return res
}

console.log(
  findDiagonalOrder([
    [1, 2, 3],
    [4, 5, 6],
    [7, 8, 9],
  ])
)
// [1,2,4,7,5,3,6,8,9]
export {}
