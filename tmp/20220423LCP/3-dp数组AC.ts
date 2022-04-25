const POS: number[][] = []
for (const r of [0, 1, 2]) {
  for (const c of [0, 1, 2]) {
    POS.push([r, c])
  }
}

// 注意js的cache很容易TLE
// 要优化就要把dfs优化成dp 这样常数会小一些
function getMaximumNumber(moles: number[][]): number {
  const times = new Set<number>([0])
  const record = new Map<number, Set<number>>()
  for (const [t, r, c] of moles) {
    times.add(t)
    !record.has(t) && record.set(t, new Set())
    record.get(t)!.add(r * 3 + c)
  }

  const allTimes = [...times].sort((a, b) => a - b)
  const n = allTimes.length

  // 也可以Int32Array
  const dp = Array.from({ length: n }, () => [
    [0, 0, 0],
    [0, 0, 0],
    [0, 0, 0],
  ])

  // i=0 情况
  dp[0][1][1] = Number(record.has(allTimes[0]) && record.get(allTimes[0])!.has(4))

  for (let i = 1; i < n; i++) {
    const curTime = allTimes[i]
    const preTime = allTimes[i - 1]
    const diff = curTime - preTime
    for (const [r, c] of POS) {
      for (const [preR, preC] of POS) {
        const cur = Number(record.has(curTime) && record.get(curTime)!.has(r * 3 + c))
        if (Math.abs(r - preR) + Math.abs(c - preC) <= diff) {
          dp[i][r][c] = Math.max(dp[i][r][c], dp[i - 1][preR][preC] + cur)
        }
      }
    }
  }

  return Math.max(...dp[n - 1].flat())
}

// console.log(
//   getMaximumNumber([
//     [1, 1, 0],
//     [2, 0, 1],
//     [4, 2, 2],
//   ])
// )

// console.log(
//   getMaximumNumber([
//     [2, 0, 2],
//     [6, 2, 0],
//     [4, 1, 0],
//     [2, 2, 2],
//     [3, 0, 2],
//   ])
// )
// 4
// console.log(
//   getMaximumNumber([
//     [4, 2, 1],
//     [1, 1, 2],
//     [2, 1, 0],
//     [3, 2, 1],
//     [4, 0, 0],
//     [0, 1, 1],
//   ])
// )

// 1
console.log(
  getMaximumNumber([
    [0, 0, 0],
    [1, 1, 0],
    [0, 2, 0],
    [1, 0, 1],
    [1, 2, 1],
  ])
)
export {}
// /**
//  * @param {number[][]} moles
//  * @return {number}
//  */
//  var getMaximumNumber = function(moles) {
//   let map = new Map()
//   for (let [t, x, y] of moles) {
//       if (!map.has(t)) map.set(t, mkCell())
//       map.get(t)[x][y] = 1
//   }
//   if (!map.has(0)) {
//       map.set(0, mkCell())
//   }
//   let sorted = [...map.entries()].sort((l, r) => l[0] - r[0])
//   let N = sorted.length
//   sorted.push([Infinity, mkCell()])
//   // console.log(sorted)

//   let dp = Array(N + 1).fill(0).map(mkCell)
//   // console.log(dp)

//   for (let i = N - 1; i >= 0; i--) {
//       for (let [x0, y0] of forItem()) {
//           let [t0, cell0] = sorted[i]
//           let value = cell0[x0][y0]

//           for (let [x1, y1] of forItem()) {
//               let [t1, cell1] = sorted[i + 1]
//               let dist = Math.abs(x0 - x1) + Math.abs(y0 - y1)
//               if (t1 - t0 >= dist) {
//                   dp[i][x0][y0] = Math.max(
//                       dp[i][x0][y0],
//                       dp[i + 1][x1][y1] + value
//                   )
//               }
//           }
//       }
//       // console.log(`dp[${i}] = ${dp[i]}`)
//   }
//   return dp[0][1][1]
//   function* forItem() {
//       for (let x = 0; x < 3; x++) {
//           for (let y = 0; y < 3; y++) {
//               yield [x, y]
//           }
//       }
//   }
//   // dp[i][x][y] =
//   function mkCell() {
//       return Array(3).fill(0).map(x => Array(3).fill(0))
//   }
// };
