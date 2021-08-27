// 请不要使用除法，且在 O(n) 时间复杂度内完成此题。
function productExceptSelf(nums: number[]): number[] {
  const len = nums.length
  const preLeft = Array<number>(len).fill(1)
  const preRight = Array<number>(len).fill(1)
  preLeft.reduce((pre, _, index, arr) => (arr[index] = pre * nums[index - 1]))
  preRight.reduceRight((pre, _, index, arr) => (arr[index] = pre * nums[index + 1]))

  return Array.from({ length: len }, (_, i) => preLeft[i] * preRight[i])
}

console.log(productExceptSelf([1, 2, 3, 4]))
// 输出: [24,12,8,6]
