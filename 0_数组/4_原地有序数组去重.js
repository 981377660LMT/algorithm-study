/**
 * @param {number[]} nums
 * @return {number}
 */
// var removeDuplicates = function (nums) {
//   let count = 0
//   for (let i = 0; i < nums.length; i++) {
//     if (nums[i] === nums[i - 1]) {
//       nums.splice(i, 1)
//       i--
//     } else {
//       count++
//     }
//   }
//   return count
// }
var removeDuplicates = function (nums) {
  // 出现了几个不同的数
  let count = 0
  for (let i = 0; i < nums.length; i++) {
    if (nums[i] !== nums[count]) {
      count++
      nums[count] = nums[i]
    }
  }
  return count + 1
}

console.log(removeDuplicates([0, 0, 1, 1, 1, 2, 2, 3, 3, 4]))
