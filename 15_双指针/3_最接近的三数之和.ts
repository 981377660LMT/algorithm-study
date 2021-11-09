// 找出 nums 中的三个整数，使得它们的和与 target 最接近。返回这三个数的和
// 假定每组输入只存在唯一答案。
// 不需要查找表，因为不需要特定的值
// 定一移二
const threeSumClosest = (nums: number[], target: number): number => {
  // 与target的差值
  let closestSum = Infinity
  nums.sort((a, b) => a - b)

  for (let l = 0; l < nums.length - 2; l++) {
    let i = l + 1
    let r = nums.length - 1

    while (i < r) {
      const sum = nums[l] + nums[i] + nums[r]

      if (Math.abs(sum - target) < Math.abs(closestSum - target)) closestSum = sum
      if (sum === target) {
        return target
      } else if (sum > target) {
        r--
      } else {
        i++
      }
    }

    // while (nums[l] === nums[l + 1]) l++
  }

  return closestSum
}

console.log(threeSumClosest([-1, 2, 1, -4], 1))

export {}
