// 这里不用位运算O(n)而是二分的方法O(logn) 因为数组有序
/**
 * @param {number[]} nums   长度为奇数
 * @return {number}
 * 当mid为偶数时，mid两边的数字个数为偶数个
 * 统一mid为奇数和偶数的情况，当mid为奇数的时候，将mid左移一位
 * mid与mid+1的值相等的话在mid的左边，不相等的话在mid的右边
 * 不断二分查找剩余元素个数为奇数的一边
 */
const singleNonDuplicate = function (nums: number[]): number {
  let l = 0
  let r = nums.length - 1

  // 因此当 left <= right 的时候，解空间都不为空，此时我们都需要继续搜索
  while (l <= r) {
    let mid = (l + r) >> 1
    if (mid % 2 === 1) mid-- // 仅对偶数索引进行二分搜索
    if (nums[mid] === nums[mid + 1]) l = mid + 2
    else r = mid - 1
  }

  return nums[l]
}

console.log(singleNonDuplicate([1, 1, 2, 3, 3, 4, 4, 8, 8]))
