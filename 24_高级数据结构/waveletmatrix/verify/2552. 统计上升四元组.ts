import { WaveletMatrix } from '../WaveletMatrix'

// eslint-disable-next-line no-inner-declarations
function countQuadruplets(nums: number[]) {
  const W = new WaveletMatrix(new Uint32Array(nums))
  let res = 0
  for (let j = 1; j < nums.length - 2; j++) {
    for (let k = j + 1; k < nums.length - 1; k++) {
      if (nums[k] < nums[j]) {
        const left = W.countRange(0, j, 0, nums[k])
        const right = W.countRange(k + 1, nums.length, nums[j] + 1, 1e9)
        res += left * right
      }
    }
  }
  return res
}
