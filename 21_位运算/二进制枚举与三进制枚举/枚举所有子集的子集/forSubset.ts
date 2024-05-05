/* eslint-disable space-in-parens */

function foo() {
  // forSubset枚举某个状态的所有子集(子集的子集)
  const state = 0b1101
  for (let g1 = state; ~g1; g1 = g1 === 0 ? -1 : (g1 - 1) & state) {
    if (g1 === state || g1 === 0) continue
    const g2 = state ^ g1
    console.log(g1.toString(2), g2.toString(2))
  }
}

/**
 * 升序枚举state所有子集的子集.
 * 0b1101 -> 0,1,4,5,8,9,12,13.
 */
function enumerateSubsetOfStateAscending(
  state: number,
  callback: (subset: number) => void | boolean
): void {
  for (let x = 0; ; x = (x - state) & state) {
    if (callback(x)) return
    if (x === state) break
  }
}

/**
 * 降序枚举state所有子集的子集.
 * 0b1101 -> 13,12,9,8,5,4,1,0.
 */
function enumerateSubsetOfStateDescending(
  state: number,
  callback: (subset: number) => void | boolean
): void {
  for (let x = state; ; x = (x - 1) & state) {
    if (callback(x)) return
    if (x === 0) break
  }
}

/**
 * 升序枚举state的所有超集.
 * 0b1101 -> 13,15.
 */
function enumerateSupersetOfState(
  n: number,
  state: number,
  callback: (superset: number) => void | boolean
): void {
  const upper = 1 << n
  for (let x = state; x < upper; x = (x + 1) | state) {
    if (callback(x)) return
  }
}

/**
 * 遍历n个元素的集合中大小为k的子集(combinations).
 * 一共有C(n,k)个子集.
 * C(4,2) -> 3,5,6,9,10,12.
 */
function enumerateSubsetOfSizeK(
  n: number,
  k: number,
  callback: (subset: number) => void | boolean
): void {
  if (k <= 0 || k > n) return
  const upper = 1 << n
  for (let x = (1 << k) - 1; x < upper; ) {
    if (callback(x)) return
    const t = x | (x - 1)
    // nextCombination (gosper hack)
    x = (t + 1) | (((~t & -~t) - 1) >>> (32 - Math.clz32(x & -x)))
  }
}

if (require.main === module) {
  enumerateSubsetOfStateAscending(0b1101, subset => {
    console.log(subset.toString(2))
  })
  enumerateSubsetOfStateDescending(0b1101, subset => {
    console.log(subset.toString(2))
  })
  enumerateSupersetOfState(4, 0b1101, superset => {
    console.log(superset.toString(2))
  })
  enumerateSubsetOfSizeK(4, 2, subset => {
    console.log(subset.toString(2), '999')
  })
  console.log('ok')
}

export {
  enumerateSubsetOfStateAscending,
  enumerateSubsetOfStateDescending,
  enumerateSubsetOfSizeK,
  enumerateSupersetOfState
}
