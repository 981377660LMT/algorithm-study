// 2～10为数字本身，A为1，J为11，Q为12，K为13，而大、小王为 0 ，
// 可以看成任意数字。

// 核心思路：数组中max-min<5且无重复元素
function isStraight(nums: number[]): boolean {
  nums.sort((a, b) => a - b)
  let joker = 0
  for (let i = 0; i < 4; i++) {
    // 统计大小王
    if (nums[i] === 0) joker++
    // 1.无重复
    else if (nums[i] === nums[i + 1]) return false
  }
  // 2. 最大减最小小于5
  return nums[4] - nums[joker] < 5
}

console.log(isStraight([0, 0, 1, 2, 5]))
// 1 2 0 0 5  ok
