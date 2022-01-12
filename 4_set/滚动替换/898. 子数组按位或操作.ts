/**
 * @param {number[]} arr
 * @return {number}
 * @summary
 * 子数组具有连续性:考虑使用前缀
 */
const subarrayBitwiseORs = (arr: number[]): number => {
  let pre = new Set<number>([0])
  const res = new Set<number>()

  for (const num of arr) {
    const cur = new Set<number>()
    for (const p of pre) {
      cur.add(num | p)
      cur.add(num)
    }

    for (const c of cur) {
      res.add(c)
    }

    pre = cur
  }

  return res.size
}

console.log(subarrayBitwiseORs([1, 1, 2]))
// 可能的子数组为 [1]，[1]，[2]，[1, 1]，[1, 2]，[1, 1, 2]。
// 产生的结果为 1，1，2，1，3，3 。

export {}
