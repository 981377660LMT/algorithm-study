zip(['a', 'b'], [1, 2], [true, false]) // [['a', 1, true], ['b', 2, false]]
console.log(zip(['a'], [1, 2], [true, false])) // [['a', 1, true], [undefined, 2, false]]

function zip<T>(...arrs: any[][]) {
  const maxLength = Math.max(...arrs.map(arr => arr.length))
  return Array.from({ length: maxLength }, (_, col) =>
    Array.from({ length: arrs.length }, (_, row) => arrs[row][col])
  )
}
export {}
