const res = [1, 2]
console.log(new Set([res, res]))
// Set(1) { [ 1, 2 ] }
console.log(
  new Set([
    [1, 2],
    [1, 2]
  ])
)
// Set(2) { [ 1, 2 ], [ 1, 2 ] }
export {}
