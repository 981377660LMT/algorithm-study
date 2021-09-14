/**
 * @param {number[]} nums
 * @param {number} F
 * @param {number} S
 * @return {number}
 * 给出非负整数数组 A ，返回两个非重叠（连续）子数组中元素的最大和，子数组的长度分别为 L 和 M。
 * @summary
 *  分别以L在前M在后，M在前L在后两种方式，滑动窗口，每次移动一格。
 */
var maxSumTwoNoOverlap = function (nums: number[], F: number, S: number): number {
  const pre = nums.slice()
  pre.reduce((pre, _, index, array) => (array[index] += pre))

  console.log(pre)
  let res = pre[F + S - 1]
  let firstMax = pre[F - 1]
  let secondMax = pre[S - 1]

  // i代表当前位于右边的数组的末尾索引
  for (let i = F + S; i < nums.length; i++) {
    //后面留一段给M， 前面 L 的最大和
    firstMax = Math.max(firstMax, pre[i - S] - pre[i - F - S])
    //后面留一段给L， 前面 M 的最大和
    secondMax = Math.max(secondMax, pre[i - F] - pre[i - F - S])
    // 前面是 L + 当前的 M
    // 前面是 M + 当前的 L
    res = Math.max(res, Math.max(firstMax + pre[i] - pre[i - S], secondMax + pre[i] - pre[i - F]))
  }

  return res
}

console.log(maxSumTwoNoOverlap([0, 6, 5, 2, 2, 5, 1, 9, 4], 1, 2))

export {}
