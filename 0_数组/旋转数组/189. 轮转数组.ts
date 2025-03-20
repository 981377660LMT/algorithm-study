import { rotateRight } from '../原地操作/rotate'
/**
 Do not return anything, modify nums in-place instead.
 */
function rotate(nums: number[], k: number): void {
  rotateRight(nums, 0, nums.length, k)
}
