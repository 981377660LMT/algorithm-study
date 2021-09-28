/**
 * @param {number[]} nums
 * @param {number} lower
 * @param {number} upper
 * @return {string[]}
 * @summary
 * 1.空的时候处理
 * 2.lower和upper特殊处理
 * 3.中间处理 cur 与 next
 */
const findMissingRanges = function (nums: number[], lower: number, upper: number): string[] {
  if (!nums.length) return lower === upper ? [`${lower}`] : [`${lower}->${upper}`]
  const res: number[][] = []

  const first = nums[0]
  if (first - lower === 1) res.push([lower])
  else if (first - lower >= 2) res.push([lower, first - 1])

  for (let i = 0; i < nums.length - 1; i++) {
    const cur = nums[i]
    const next = nums[i + 1]
    const diff = next - cur
    if (diff <= 1) continue
    else if (diff === 2) res.push([cur + 1])
    else res.push([cur + 1, next - 1])
  }

  const last = nums[nums.length - 1]
  if (upper - last === 1) res.push([upper])
  else if (upper - last >= 2) res.push([last + 1, upper])

  return res.map(item => (item.length === 1 ? item[0].toString() : `${item[0]}->${item[1]}`))
}

console.log(findMissingRanges([0, 1, 3, 50, 75], 0, 99))
console.log(findMissingRanges([-1], -2, -1))
// 输出: ["2", "4->49", "51->74", "76->99"]
