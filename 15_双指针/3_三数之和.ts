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
  for (let p1 = 0; p1 < nums.length - 2; p1++) {
    // 定一(最左边的数要不超过0，否则三数之和必不为0)
    // 加速循环
    if (nums[p1] > target) break

    let p2 = p1 + 1
    let p3 = nums.length - 1
    while (p2 < p3) {
      const curSum = nums[p1] + nums[p2] + nums[p3]
      if (curSum === target) {
        res.push([nums[p1], nums[p2], nums[p3]])

        // 这里需要去除重复!
        while (nums[p2] === nums[p2 + 1]) p2++
        while (nums[p3] === nums[p3 - 1]) p3--

        p2++
        p3--
      } else if (curSum < target) {
        p2++
      } else {
        p3--
      }
    }
    // "定一"的元素相等时跳过本次循环避免重复
    while (nums[p1] === nums[p1 + 1]) p1++
  }

  return res
}

console.log(threeSum([-1, 0, 1, 2, -1, -4]))
console.log(threeSum([-2, 0, 0, 2, 2]))
// 输出：[[-1,-1,2],[-1,0,1]]

export {}
