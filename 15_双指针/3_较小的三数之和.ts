// 寻找能够使条件 nums[i] + nums[j] + nums[k] < target 成立的三元组  i, j, k 个数
function threeSumSmaller(nums: number[], target: number): number {
  if (nums.length < 3) return 0

  nums.sort((a, b) => a - b)
  let res = 0

  for (let p1 = 0; p1 < nums.length - 2; p1++) {
    // 定一移二
    let p2 = p1 + 1
    let p3 = nums.length - 1
    while (p2 < p3) {
      const curSum = nums[p1] + nums[p2] + nums[p3]
      if (curSum < target) {
        res += p3 - p2 // 右边全移过来
        p2++
      } else {
        p3--
      }
    }
  }

  return res
}

console.log(threeSumSmaller([-2, 0, 1, 3], 2))
