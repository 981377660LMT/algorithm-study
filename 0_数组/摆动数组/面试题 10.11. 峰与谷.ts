/**
 Do not return anything, modify nums in-place instead.
 */
function wiggleSort(nums: number[]): void {
  nums.sort((a, b) => a - b)
  let mid = nums.length >> 1
  const small = nums.slice(0, mid)
  const big = nums.slice(mid)
  console.log(small, big)
  for (let i = 0; i < nums.length; i++) {
    nums[i] = i & 1 ? small.pop()! : big.pop()!
  }
}

const testArr = [5, 3, 1, 2, 3, 2, 44, 23]
wiggleSort(testArr)
console.log(testArr)
