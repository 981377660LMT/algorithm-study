// 所以当flag从1变到2，或者从2变到1，就是波峰波谷
function countHillValley(nums: number[]): number {
  let res = 0
  let pre = 0

  for (let i = 1; i < nums.length; i++) {
    if (nums[i] > nums[i - 1]) {
      if (pre === 1) res++
      pre = 2
    } else if (nums[i] < nums[i - 1]) {
      if (pre === 2) res++
      pre = 1
    }
  }

  return res
}

console.log(countHillValley([2, 4, 1, 1, 6, 5]))
