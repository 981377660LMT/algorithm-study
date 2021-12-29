/**
 * @param {number} num 给定数字的范围是 [0, 108] 8位最多28种 暴力即可O(n^2)
 * @return {number}
 * 给定一个非负整数，你至多可以交换一次数字中的任意两位。返回你能得到的最大值。
 * @summary 让第一个高位的小数和后面的尽量大的大数交换 大数要尽可能后
 * 记录每个数出现的最后索引即可
 * @link https://leetcode-cn.com/problems/maximum-swap/solution/dan-diao-zhan-yu-tan-xin-gui-lu-by-user5707f/
 */
function maximumSwap(num: number): number {
  const nums = num.toString().split('').map(Number)
  const lastIndex = new Map<number, number>()
  nums.forEach((val, index) => lastIndex.set(val, index))

  for (const [index, val] of nums.entries()) {
    for (let cand = 9; cand > val; cand--) {
      // 后面可交换的数
      if ((lastIndex.get(cand) || -1) > index) {
        const swapIndex = lastIndex.get(cand)!
        ;[nums[index], nums[swapIndex]] = [nums[swapIndex], nums[index]]
        return Number(nums.join(''))
      }
    }
  }

  return num
}

console.log(maximumSwap(2736))

export default 1
