/* eslint-disable consistent-return */
/* eslint-disable no-loop-func */
/* eslint-disable no-inner-declarations */

import { stringArrayToNumber } from '../字符串与数字/stringArrayToNumber'

/**
 * 遍历数组所有的分割方案.
 * @param f 返回 true 时停止遍历.
 */
function enumerateSplits<T>(arr: T[], f: (groups: T[][]) => boolean | void): void {
  if (!arr.length) return
  const n = arr.length
  const dfs = (pos: number, groups: T[][]): boolean => {
    if (pos === n) {
      return !!f(groups)
    }
    for (let i = pos; i < n; i++) {
      groups.push(arr.slice(pos, i + 1))
      if (dfs(i + 1, groups)) return true
      groups.pop()
    }
    return false
  }
  dfs(0, [])
}

/**
 * 遍历数组所有的分割方案.
 * @deprecated
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

export { enumerateSplits }

if (require.main === module) {
  enumerateSplits([1, 2, 3], console.log)

  // https://leetcode.cn/problems/find-the-punishment-number-of-an-integer/
  function punishmentNumber(n: number): number {
    let res = 0
    for (let i = 1; i <= n; i++) {
      const sb = String(i * i).split('')

      enumerateSplits(sb, groups => {
        let sum = 0
        groups.forEach(g => {
          sum += stringArrayToNumber(g)
        })

        if (sum === i) {
          res += i * i
          return true
        }
      })
    }

    return res
  }

  console.log(punishmentNumber(10))

  console.time('enumerateSplits')
  enumerateSplits(
    Array.from({ length: 25 }, (_, i) => i),
    () => {}
  )
  console.timeEnd('enumerateSplits')

  console.time('enumerateSplitsLegacy')
  enumerateSplitsLegacy(
    Array.from({ length: 25 }, (_, i) => i),
    () => {}
  )
  console.timeEnd('enumerateSplitsLegacy')
}
