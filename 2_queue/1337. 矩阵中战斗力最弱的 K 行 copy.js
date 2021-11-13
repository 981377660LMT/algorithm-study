const kWeakestRows = function (mat, k) {
  return (
    mat
      //turn the array into [sum of soilders,index] form
      .map((e, i) => [e.reduce((acc, cur) => acc + cur, 0), i])
      //sort the array: if the number of soilders is equal then sort with the index of the row
      .sort((a, b) => a[0] - b[0] || a[1] - b[1])
      //take of the row index of the sorted result
      .map(x => x[1])
      //slice the result according to k number
      .slice(0, k)
  )
}

console.log(
  kWeakestRows(
    [
      [1, 1, 0, 0, 0],
      [1, 1, 1, 1, 0],
      [1, 0, 0, 0, 0],
      [1, 1, 0, 0, 0],
      [1, 1, 1, 1, 1],
    ],
    3
  )
)
