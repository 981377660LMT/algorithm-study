/**
 * @param {number[]} nums
 * @return {number}
 * 三角形的个数
 */
function triangleNumber(nums) {
  nums.sort((a, b) => a - b)

  let res = 0
  for (let p3 = nums.length - 1; p3 >= 2; p3--) {
    let p1 = 0
    let p2 = p3 - 1
    while (nums[p1] === 0) p1++
    while (p1 < p2) {
      if (nums[p1] + nums[p2] > nums[p3]) {
        res += p2 - p1
        p2--
      } else {
        p1++
      }
    }
  }

  return res
}
