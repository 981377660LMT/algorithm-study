/**
 * 将 num 拆分成 k 和 k+1 的和，使得拆分的个数最(多/少).
 * @param num 正整数.
 * @param k 正整数.
 * @param minimize 是否使得拆分的个数最少. 默认为最少(true).
 * @returns [count1, count2, ok] count1和count2分别是拆分成k和k+1的个数，ok表示是否可以拆分.
 */
function splitToKAndKPlusOne(
  num: number,
  k: number,
  minimize = true
): [count1: number, count2: number, ok: boolean] {
  if (minimize) {
    const count2 = Math.ceil(num / (k + 1))
    const diff = (k + 1) * count2 - num
    if (diff > count2) return [0, 0, false]
    return [diff, count2 - diff, true]
  }

  const count1 = Math.floor(num / k)
  const diff = num - k * count1
  if (diff > count1) return [0, 0, false]
  return [count1 - diff, diff, true]
}

export { splitToKAndKPlusOne }

if (require.main === module) {
  console.log(splitToKAndKPlusOne(12, 3))
  console.log(splitToKAndKPlusOne(12, 3, false))
  console.log(splitToKAndKPlusOne(5, 2))
  console.log(splitToKAndKPlusOne(5, 2, false))
  console.log(splitToKAndKPlusOne(16, 3))
  console.log(splitToKAndKPlusOne(16, 3, false))

  // eslint-disable-next-line no-inner-declarations
  function checkWithBruteForce(
    num: number,
    k: number,
    minimize = true
  ): [count1: number, count2: number, ok: boolean] {
    let res = [0, 0, false]
    for (let count1 = 0; count1 <= num; count1++) {
      for (let count2 = 0; count2 <= num; count2++) {
        if (count1 * k + count2 * (k + 1) === num) {
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

  for (let num = 0; num < 100; num++) {
    for (let k = 1; k < 100; k++) {
      for (let minimize = 0; minimize < 2; minimize++) {
        const res1 = splitToKAndKPlusOne(num, k, Boolean(minimize))
        const res2 = checkWithBruteForce(num, k, Boolean(minimize))
        if (res1[2] !== res2[2] || res1[0] !== res2[0] || res1[1] !== res2[1]) {
          console.error(num, k, minimize)
          console.error(res1)
          console.error(res2)
          throw new Error()
        }
      }
    }
  }

  console.log('pass!')

  // 2870. 使数组为空的最少操作次数
  // https://leetcode.cn/problems/minimum-number-of-operations-to-make-array-empty/
  // eslint-disable-next-line no-inner-declarations
  function minOperations(nums: number[]): number {
    const counter = new Map<number, number>()
    nums.forEach(v => counter.set(v, (counter.get(v) || 0) + 1))
    let res = 0
    for (const count of counter.values()) {
      const [count1, count2, ok] = splitToKAndKPlusOne(count, 2, true)
      if (!ok) return -1
      res += count1 + count2
    }
    return res
  }
}
