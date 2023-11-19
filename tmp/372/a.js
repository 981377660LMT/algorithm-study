/**
 * @param {number[]} heights
 * @param {number[][]} queries
 * @return {number[]}
 */
var leftmostBuildingQueries = function (h, q) {
  const n = h.length
  const m = q.length
  const res = Array(m)
  for (let i = 0; i < m; i++) {
    res[i] = -1
    const q1 = q[i][0]
    const q2 = q[i][1]
    if (q1 == q2) {
      res[i] = q1
    } else if (q1 < q2 && h[q1] < h[q2]) {
      res[i] = q2
    } else if (q2 < q1 && h[q2] < h[q1]) {
      res[i] = q1
    } else {
      let ma = q1
      if (ma < q2) {
        ma = q2
      }
      for (let j = ma + 1; j < n; j++) {
        if (h[j] > h[q1] && h[j] > h[q2]) {
          res[i] = j
          break
        }
      }
    }
  }
  return res
}
