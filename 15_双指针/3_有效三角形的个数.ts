/**
 * @param {number[]} nums
 * @return {number}
 * @description 固定最长的一条边比较方便 定一移二
 */
const triangleNumber = (nums: number[]): number => {
  nums.sort((a, b) => a - b)
  let res = 0

  for (let p3 = nums.length - 1; p3 >= 2; p3--) {
    let p1 = 0
    let p2 = p3 - 1
    while (nums[p1] === 0) p1++
    while (p1 < p2) {
      if (nums[p1] + nums[p2] > nums[p3]) {
        res += p2 - p1 // 左边全部移过去
        p2--
      } else {
        p1++
      }
    }
  }

  return res
}

console.log(triangleNumber([2, 2, 3, 4]))
console.log(triangleNumber([4, 2, 3, 4]))
// 输出: 3
// 解释:
// 有效的组合是:
// 2,3,4 (使用第一个 2)
// 2,3,4 (使用第二个 2)
// 2,2,3

export {}
