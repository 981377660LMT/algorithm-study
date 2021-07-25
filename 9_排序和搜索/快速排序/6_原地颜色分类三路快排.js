// 仅使用常数空间的一趟扫描算法吗？
// /**
//  * @param {number[]} nums
//  * @return {void} Do not return anything, modify nums in-place instead.
//  * @description 普通的计数排序,整个数组扫描了两便
//  */
// var sortColors = function (nums) {
//   let zero = 0
//   let one = 0
//   let two = 0

//   for (let i = 0; i < nums.length; i++) {
//     switch (nums[i]) {
//       case 0:
//         zero++
//         break
//       case 1:
//         one++
//         break
//       case 2:
//         two++
//         break
//       default:
//         break
//     }
//   }

//   for (let i = 0; i < zero; i++) {
//     nums[i] = 0
//   }
//   for (let i = zero; i < zero + one; i++) {
//     nums[i] = 1
//   }
//   for (let i = zero + one; i < nums.length; i++) {
//     nums[i] = 2
//   }

//   return nums
// }
/**
 * @param {number[]} nums
 * @return {void} Do not return anything, modify nums in-place instead.
 * @description 三路快排一遍扫描 找一个基准v，通过双向指针，把<v的值放在左边，>v的值放在右边，等于v的放中间。
 */
var sortColors = function (nums) {
  let l = 0
  let i = 0
  let r = nums.length - 1

  const swap = (i, j) => {
    ;[nums[i], nums[j]] = [nums[j], nums[i]]
  }

  // 注意结束
  while (i <= r) {
    switch (nums[i]) {
      case 0:
        swap(i, l)
        i++
        l++
        break
      case 1:
        i++
        break
      case 2:
        swap(i, r)
        r--
        break
      default:
        break
    }
  }

  return nums
}
console.log(sortColors([2, 0, 1, 0, 1, 2, 0, 1]))
// 输出：[0,0,1,1,2,2]
