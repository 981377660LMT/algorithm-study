import assert from 'assert'
import { SparseTable } from './SparseTable'

function closestToTarget(arr: number[], target: number): number {
  const n = arr.length
  const st = new SparseTable(arr, (a, b) => a & b)

  let res = Math.abs(arr[0] - target)
  for (let start = 0; start < n; start++) {
    let [left, right] = [start, n - 1]
    while (left <= right) {
      const mid = Math.floor((left + right) / 2)
      const diff = st.query(start, mid) - target
      res = Math.min(res, Math.abs(diff))
      if (diff === 0) return 0
      if (diff > 0) left = mid + 1
      else right = mid - 1
    }
  }

  return res
}

if (require.main === module) {
  assert.strictEqual(closestToTarget([9, 12, 3, 7, 15], 5), 2)
  assert.strictEqual(closestToTarget([70, 15, 21, 96], 4), 0)
}
