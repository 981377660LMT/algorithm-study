// 小技巧，搜索前对字符串变换
// 1764. 通过连接另一个数组的子数组得到一个数组
// function check(nums: number[]): boolean {
//   const str1 = nums.map(num => `#${num}#`).join('')
//   const str2 = nums
//     .slice()
//     .sort((a, b) => a - b)
//     .map(num => `#${num}#`)
//     .join('')

//   return isFlipedString(str1, str2)
// }

// 至多存在一个下降点
function sum2(nums: number[]): boolean {
  let down = 0

  for (let i = 0; i < nums.length; i++) {
    if (nums[i] > nums[(i + 1) % nums.length]) down++
  }

  return down <= 1
}

console.log(sum2([3, 4, 5, 1, 2]))
