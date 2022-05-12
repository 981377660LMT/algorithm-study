/**
 * @param {number[]} nums
 * @return {number[]}
 * @description 给定一个循环数组,你应该循环地搜索它的下一个更大的数。如果不存在，则输出 -1。
 */
function nextGreaterElements(nums) {
  const n = nums.length
  const res = Array(n).fill(-1)
  const stack = []

  // i全部模
  for (let i = 0; i < 2 * n; i++) {
    // 下一个更大
    while (stack.length && nums[stack[stack.length - 1]] < nums[i % n]) {
      const tmp = stack.pop()
      res[tmp] = nums[i % n]
    }
    stack.push(i % n)
  }

  return res
}

console.log(nextGreaterElements([5, 4, 3, 2, 1]))
// [-1,5,5,5,5]
