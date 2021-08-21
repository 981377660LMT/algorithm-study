/**
 * @param {number[]} nums
 * @return {void} Do not return anything, modify nums in-place instead.
 * 如果不存在下一个更大的排列，则将数字重新排列成最小的排列（即升序排列）。
 */
const nextPermutation = function (nums: number[]): void {
  let mid = 0
  let mono = true
  for (let index = 0; index < array.length; index++) {
    const element = array[index]
  }

  if (mono) nums.reverse()
}

console.log(nextPermutation([1, 2, 3]))
// 输出：[1,3,2]
console.log(nextPermutation([3, 2, 1]))
// 输出：[1,2,3]
export default 1
