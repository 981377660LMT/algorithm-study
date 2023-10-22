/**
 * 将 num 拆分成 a 和 b 的和，使得拆分的个数最(多/少).
 * @param num 正整数.
 * @param a 正整数.
 * @param b 正整数.
 * @param minimize 是否使得拆分的个数最少. 默认为最少(true).
 * @returns [countA, countB, ok] countA和countB分别是拆分成a和b的个数，ok表示是否可以拆分.
 */
function splitToAAndB(
  num: number,
  a: number,
  b: number,
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

export { splitToAAndB }
