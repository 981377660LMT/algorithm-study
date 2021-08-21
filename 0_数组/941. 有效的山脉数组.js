// 0 < i < arr.length - 1 条件下，存在 i 使得：
// arr[0] < arr[1] < ... arr[i-1] < arr[i]
// arr[i] > arr[i+1] > ... > arr[arr.length - 1]

// 左右双指针 必须要在中间相遇
/**
 * @param {number[]} arr
 * @return {boolean}
 */
var validMountainArray = function (arr) {
  let left = 0,
    right = arr.length - 1

  while (arr[left] < arr[left + 1]) left++
  while (arr[right] < arr[right - 1]) right--

  if (left !== right || left === arr.length - 1 || right === 0) {
    return false
  }

  return true
}

console.log(validMountainArray([0, 3, 2, 1]))
