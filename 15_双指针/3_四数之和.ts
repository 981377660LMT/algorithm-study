// 四指针 定二移二
// 同理n指针 定n-2移2
// 定的元素要判断是否与上一个重复 移的元素要判断与下一个是否重复
const fourSum = (nums: number[], target: number): number[][] => {
  if (nums.length < 4) return []
  nums.sort((a, b) => a - b)
  const res: number[][] = []

  for (let p1 = 0; p1 < nums.length - 3; p1++) {
    for (let p2 = p1 + 1; p2 < nums.length - 2; p2++) {
      let p3 = p2 + 1
      let p4 = nums.length - 1
      while (p3 < p4) {
        const sum = nums[p1] + nums[p2] + nums[p3] + nums[p4]
        if (sum === target) {
          res.push([nums[p1], nums[p2], nums[p3], nums[p4]])
          while (nums[p4] === nums[p4 - 1]) p4--
          while (nums[p3] === nums[p3 + 1]) p3++
          p4--
          p3++
        } else if (sum > target) {
          p4--
        } else {
          p3++
        }
      }
      while (nums[p2] === nums[p2 + 1]) p2++
    }
    while (nums[p1] === nums[p1 + 1]) p1++
  }

  return res
}

console.log(fourSum([1, 0, -1, 0, -2, 2], 0))
console.log(fourSum([2, 2, 2, 2, 2], 8))

export {}
