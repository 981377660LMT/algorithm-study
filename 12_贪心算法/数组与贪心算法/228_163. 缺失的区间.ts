/**
 * @param {number[]} nums
 * @param {number} lower
 * @param {number} upper
 * @return {string[]}
 * @summary
 * 根据数据范围要按nums遍历，不能按lower到upper
 * 添加Sentinel
 * 比较前后的diff是否大于2
 */
const findMissingRanges = function (nums: number[], lower: number, upper: number): string[] {
  // Sentinel
  nums.push(upper + 1)

  const res: string[] = []
  let pre = lower - 1

  for (const cur of nums) {
    const diff = cur - pre
    if (diff > 2) {
      res.push(`${pre + 1}->${cur - 1}`)
    } else if (diff === 2) {
      res.push(`${pre + 1}`)
    }

    pre = cur
  }

  return res
}

console.log(findMissingRanges([0, 1, 3, 50, 75], 0, 99))
console.log(findMissingRanges([-1], -2, -1))
// 输出: ["2", "4->49", "51->74", "76->99"]
