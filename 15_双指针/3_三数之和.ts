// 关键在于先排序
// 定一移二
// 去除重复
const threeSum = (nums: number[]) => {
  if (nums.length < 3) return []
  nums.sort((a, b) => a - b)
  const res: number[][] = []
  // 三数之和为0
  const target = 0

  console.log(nums)
  for (let index = 0; index < nums.length; index++) {
    const element = nums[index]
    // 定一(最左边的数要不超过0，否则三数之和必不为0)
    if (element > target) break
    // "定一"的元素相等时跳过本次循环避免重复
    if (index > 0 && nums[index] === nums[index - 1]) continue

    let leftPoint = index + 1
    let rightPoint = nums.length - 1
    while (leftPoint < rightPoint) {
      const curSum = element + nums[leftPoint] + nums[rightPoint]
      if (curSum === target) {
        res.push([element, nums[leftPoint], nums[rightPoint]])

        // 这里需要去除重复!
        while (nums[leftPoint] === nums[leftPoint + 1]) leftPoint++
        while (nums[rightPoint] === nums[rightPoint - 1]) rightPoint--

        leftPoint++
        rightPoint--
      }

      if (curSum < 0) leftPoint++
      if (curSum > 0) rightPoint--
    }
  }

  return res
}

console.log(threeSum([-1, 0, 1, 2, -1, -4]))
console.log(threeSum([-2, 0, 0, 2, 2]))
// 输出：[[-1,-1,2],[-1,0,1]]

export {}
