/**
 * @param {number[]} arr
 * @return {number}
 * @summary
 * 子数组具有连续性:考虑使用前缀
 */
const subarrayBitwiseORs = (arr: number[]): number => {
  let dp = new Set<number>([0])
  const res = new Set<number>()

  for (const num of arr) {
    const ndp = new Set<number>()
    for (const p of dp) {
      ndp.add(num | p)
      ndp.add(num)
    }

    for (const c of ndp) {
      res.add(c)
    }

    dp = ndp
  }

  return res.size
}

console.log(subarrayBitwiseORs([1, 1, 2]))
// 可能的子数组为 [1]，[1]，[2]，[1, 1]，[1, 2]，[1, 1, 2]。
// 产生的结果为 1，1，2，1，3，3 。

export {}
