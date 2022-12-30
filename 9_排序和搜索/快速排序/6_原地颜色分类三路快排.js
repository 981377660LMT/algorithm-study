// 仅使用常数空间的一趟扫描算法吗？

/**
 * @param {number[]} nums
 * @return {void} Do not return anything, modify nums in-place instead.
 * @description 三路快排一遍扫描 找一个基准v，通过双向指针，把<v的值放在左边，>v的值放在右边，等于v的放中间。
 */
function sortColors(nums) {
  let l = 0
  let m = 0
  let r = nums.length - 1

  const swap = (i, j) => {
    ;[nums[i], nums[j]] = [nums[j], nums[i]]
  }

  // 注意结束
  while (m <= r) {
    switch (nums[m]) {
      case 0:
        swap(m, l)
        m++
        l++
        break
      case 1:
        m++
        break
      case 2:
        swap(m, r)
        r-- // 此时不能移动m 因为要处理新过来的数

        break
      default:
        break
    }
  }

  return nums
}
console.log(sortColors([2, 0, 1, 0, 1, 2, 0, 1]))
// 输出：[0,0,1,1,2,2]
