/* eslint-disable no-inner-declarations */

/**
 * 遍历数组所有的分割方案.
 */
function enumerateSplits<T>(arr: T[], f: (splits: T[][]) => void): void {
  if (!arr.length) return
  const n = arr.length
  for (let state = 0; state < 1 << (n - 1); state++) {
    let preSplit = 0
    const cur: T[][] = []
    for (let i = 0; i < n - 1; i++) {
      if (state & (1 << i)) {
        cur.push(arr.slice(preSplit, i + 1))
        preSplit = i + 1
      }
    }
    cur.push(arr.slice(preSplit))
    f(cur)
  }
}

/**
 * 遍历数组所有的分割方案.
 */
function* genSplits<T>(arr: T[]): Generator<T[][]> {
  if (!arr.length) return
  const n = arr.length
  for (let state = 0; state < 1 << (n - 1); state++) {
    let preSplit = 0
    const cur: T[][] = []
    for (let i = 0; i < n - 1; i++) {
      if (state & (1 << i)) {
        cur.push(arr.slice(preSplit, i + 1))
        preSplit = i + 1
      }
    }
    cur.push(arr.slice(preSplit))
    yield cur
  }
}

export { enumerateSplits, genSplits }

if (require.main === module) {
  enumerateSplits([1, 2, 3], console.log)

  // https://leetcode.cn/problems/find-the-punishment-number-of-an-integer/
  function punishmentNumber(n: number): number {
    let res = 0
    for (let i = 1; i <= n; i++) {
      const sb = String(i * i).split('')
      for (const split of genSplits(sb)) {
        let sum = 0
        for (const s of split) {
          sum += +s.join('')
        }
        if (sum === i) {
          res += i * i
          break
        }
      }
    }
    return res
  }

  console.log(punishmentNumber(10))
}
