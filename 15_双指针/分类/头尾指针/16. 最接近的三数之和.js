/**
 * @param {number[]} nums
 * @param {number} target
 * @return {number}
 */
function threeSumClosest(nums, target) {
  const n = nums.length
  let res = Infinity

  nums.sort((a, b) => a - b)

  for (let i = 0; i < n - 2; i++) {
    let left = i + 1
    let right = n - 1

    while (left < right) {
      const sum = nums[i] + nums[left] + nums[right]

      if (Math.abs(sum - target) < Math.abs(res - target)) {
        res = sum
      }

      if (sum === target) {
        return target
      }

      if (sum > target) {
        right--
      } else {
        left++
      }
    }
  }

  return res
}

console.log(threeSumClosest([-1, 2, 1, -4], 1))
