/**
 * @param {number[]} arr
 * @return {number}
 * 很容易想到时间复杂度 O(n) 的解决方案，你可以设计一个 O(log(n)) 的解决方案吗
 */
var peakIndexInMountainArray = function (arr) {
  let [l, r] = [0, arr.length - 1]

  while (l <= r) {
    const mid = (l + r) >> 1

    if (arr[mid] > arr[mid + 1]) {
      r = mid - 1
    } else {
      l = mid + 1
    }
  }

  return l
}

console.log(peakIndexInMountainArray([24, 69, 100, 99, 79, 78, 67, 36, 26, 19]))
