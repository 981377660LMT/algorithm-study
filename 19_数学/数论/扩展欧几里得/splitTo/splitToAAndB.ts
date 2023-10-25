/* eslint-disable no-inner-declarations */
import { solveLinearEquation } from '../solveLinearEquation/solveLinearEquation'

/**
 * 将 num 拆分成 a 和 b 的和，使得拆分的个数最(多/少).
 * @param num 正整数.
 * @param a 正整数.
 * @param b 正整数.
 * @param minimize 是否使得拆分的个数最少. 默认为最少(true).
 * @returns [countA, countB, ok] countA和countB分别是拆分成a和b的个数，ok表示是否可以拆分.
 */
function splitToAAndB(num: number, a: number, b: number, minimize = true): [count1: number, count2: number, ok: boolean] {
  const { n, x1, y1, x2, y2 } = solveLinearEquation(a, b, num, true) // 允许解为0
  if (n < 0) return [0, 0, false]
  if (n > 0) {
    const res1Smaller = x1 + y1 <= x2 + y2
    return res1Smaller === minimize ? [x1, y1, true] : [x2, y2, true]
  }

  // 存在整数解但不存在正整数解，检查其中一项是否可以为0
  const modA = num % a
  const modB = num % b
  if (modA && modB) return [0, 0, false]
  if (modA) return [0, num / b, true]
  if (modB) return [num / a, 0, true]
  const div1 = num / a
  const div2 = num / b
  const res1Smaller = div1 <= div2
  return res1Smaller === minimize ? [div1, 0, true] : [0, div2, true]
}

export { splitToAAndB }

if (require.main === module) {
  // 2870. 使数组为空的最少操作次数
  // https://leetcode.cn/problems/minimum-number-of-operations-to-make-array-empty/
  function minOperations(nums: number[]): number {
    const counter = new Map<number, number>()
    nums.forEach(v => counter.set(v, (counter.get(v) || 0) + 1))
    let res = 0
    for (const count of counter.values()) {
      const [count1, count2, ok] = splitToAAndB(count, 2, 3, true)
      if (!ok) return -1
      res += count1 + count2
    }
    return res
  }

  // 2910. 合法分组的最少组数
  // https://leetcode.cn/problems/minimum-number-of-groups-to-create-a-valid-assignment/
  function minGroupsForValidAssignment(nums: number[]): number {
    const n = nums.length
    const tmpCounter = new Map<number, number>()
    nums.forEach(v => tmpCounter.set(v, (tmpCounter.get(v) || 0) + 1))
    const freq = [...tmpCounter.values()]
    const freqCounter = new Map<number, number>()
    freq.forEach(v => freqCounter.set(v, (freqCounter.get(v) || 0) + 1))

    let res = n
    for (let size = 1; size < n; size++) {
      let ok = true
      let cand = 0
      for (const value of freqCounter.keys()) {
        const [count1, count2, ok_] = splitToAAndB(value, size, size + 1, true)
        if (!ok_) {
          ok = false
          break
        }
        cand += (count1 + count2) * freqCounter.get(value)!
      }
      if (ok) res = Math.min(res, cand)
    }

    return res
  }

  console.log(splitToAAndB(12, 3, 4, false))

  function checkWithBruteForce(num: number, a: number, b: number, minimize = true): [count1: number, count2: number, ok: boolean] {
    let res = [0, 0, false]
    for (let count1 = 0; count1 <= num; count1++) {
      for (let count2 = 0; count2 <= num; count2++) {
        if (count1 * a + count2 * b === num) {
          if (minimize) {
            if (!res[2] || count1 + count2 < res[0] + res[1]) res = [count1, count2, true]
          } else {
            if (!res[2] || count1 + count2 > res[0] + res[1]) res = [count1, count2, true]
          }
        }
      }
    }
    return res
  }

  for (let num = 1; num < 100; num++) {
    for (let a = 1; a < 100; a++) {
      for (let b = 1; b < 100; b++) {
        if (a === b) continue
        for (let minimize = 0; minimize < 2; minimize++) {
          const res1 = splitToAAndB(num, a, b, Boolean(minimize))
          const res2 = checkWithBruteForce(num, a, b, Boolean(minimize))
          if (res1[2] !== res2[2] || res1[0] !== res2[0] || res1[1] !== res2[1]) {
            console.error(num, a, b, minimize)
            console.error(res1)
            console.error(res2)
            throw new Error()
          }
        }
      }
    }
  }

  console.log('pass')
}
