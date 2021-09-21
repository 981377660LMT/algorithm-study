type Count = number
type Match = number

/**
 * @param {number[]} nums  nums.length <= 100000
 * @param {number} target
 * @return {number[][]}
 * 设计一个算法，找出数组中两数之和为指定值的所有整数对。一个数只能属于一个数对。
 * @summary
 * 无序数组 map
 * 有序数组 双指针
 */
const pairSums = function (nums: number[], target: number): number[][] {
  const record = new Map<Match, Count>()
  const res: number[][] = []

  for (let i = 0; i < nums.length; i++) {
    const num = nums[i]
    if (record.has(num)) {
      res.push([num, target - num])
      const count = record.get(num)!
      if (count === 1) record.delete(num)
      else record.set(num, record.get(num)! - 1)
    } else {
      record.set(target - num, (record.get(target - num) || 0) + 1)
    }
  }

  return res
}

const pairSums2 = function (nums: number[], target: number): number[][] {
  nums.sort((a, b) => a - b)
  const res: number[][] = []
  let l = 0
  let r = nums.length - 1

  while (l < r) {
    const sum = nums[l] + nums[r]
    if (sum === target) {
      res.push([nums[l], nums[r]])
      l++
      r--
    } else if (sum > target) r--
    else l++
  }

  return res
}

console.log(pairSums([5, 6, 5, 6], 11))
