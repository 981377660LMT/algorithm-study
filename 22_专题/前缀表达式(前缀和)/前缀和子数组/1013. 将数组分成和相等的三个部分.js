/**
 * @param {number[]} arr  3 <= arr.length <= 5 * 104
 * @return {boolean}
 * 只有可以将其划分为三个和相等的 非空 部分时才返回 true，否则返回 false
 */
const canThreePartsEqualSum = function (arr) {
  const total = arr.reduce((pre, cur) => pre + cur, 0)
  if (total % 3 !== 0) return false

  const target = total / 3
  let hit = 0
  let sum = 0

  // 最后一个不遍历是保证存在子数组
  for (let i = 0; i < arr.length - 1; i++) {
    sum += arr[i]
    if (sum === target) {
      sum = 0
      hit++
    }
    if (hit === 2) return true
  }

  return false
}

console.log(canThreePartsEqualSum([0, 2, 1, -6, 6, -7, 9, 1, 2, 0, 1]))
// 输出：true
// 解释：0 + 2 + 1 = -6 + 6 - 7 + 9 + 1 = 2 + 0 + 1
