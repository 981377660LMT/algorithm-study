export {}

function getLargestOutlier(nums: number[]): number {
  const counter: Record<number, number> = {}
  let sum = 0
  for (const num of nums) {
    counter[num] = (counter[num] || 0) + 1
    sum += num
  }

  const sortedNums = [...new Set(nums)].sort((a, b) => b - a)
  for (const num of sortedNums) {
    const remain = sum - num
    if (remain % 2 !== 0) continue
    const tmp = remain / 2
    if (counter[tmp]) {
      if (tmp === num) {
        if (counter[num] >= 2) return num
      } else {
        return num
      }
    }
  }
  return -1
}
