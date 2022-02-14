// 这里不用位运算O(n)而是二分的方法O(logn) 因为数组有序
/**
 * @param {number[]} nums   长度为奇数
 * @return {number}
 * 当mid为偶数时，mid两边的数字个数为偶数个
 * 统一mid为奇数和偶数的情况，当mid为奇数的时候，将mid左移一位
 * mid与mid+1的值相等的话在mid的左边，不相等的话在mid的右边
 * 不断二分查找剩余元素个数为奇数的一边
 */
function singleNonDuplicate(nums: number[]): number {
  let l = 0
  let r = nums.length - 1

  while (l <= r) {
    const mid = (l + r) >> 1
    // 不需要判断mid 的奇偶性，mid 和 mid⊕1 即为每次需要比较元素的两个下标
    if (nums[mid] === nums[mid ^ 1]) l = mid + 1
    else r = mid - 1
  }

  return nums[l]
}

console.log(singleNonDuplicate([1, 1, 2, 3, 3, 4, 4, 8, 8]))
console.log(singleNonDuplicate([1, 1, 2]))
