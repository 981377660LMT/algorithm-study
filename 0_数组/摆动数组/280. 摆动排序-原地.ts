/**
 Do not return anything, modify nums in-place instead.
 给你一个无序的数组 nums, 将该数字 `原地` 重排后使得 nums[0] <= nums[1] >= nums[2] <= nums[3]...。

 @summary
 因为是从左往右，左边的肯定都是经过计算和交换的。
 只需关心右侧的
 */
function wiggleSort(nums: number[]): void {
  for (let i = 0; i < nums.length - 1; i++) {
    if ((i & 1) === 0) {
      if (nums[i] > nums[i + 1]) {
        ;[nums[i], nums[i + 1]] = [nums[i + 1], nums[i]]
      }
    } else {
      if (nums[i] < nums[i + 1]) {
        ;[nums[i], nums[i + 1]] = [nums[i + 1], nums[i]]
      }
    }
  }
}

const testArray = [3, 5, 2, 1, 6, 4]
wiggleSort(testArray)
console.log(testArray)
export {}
