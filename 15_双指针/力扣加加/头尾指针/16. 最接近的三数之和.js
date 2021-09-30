/**
 * @param {number[]} nums
 * @param {number} target
 * @return {number}
 */
var threeSumClosest = function (nums, target) {
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
  }

  return closestSum
}

console.log(threeSumClosest([-1, 2, 1, -4], 1))
