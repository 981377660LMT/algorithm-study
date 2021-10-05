// // 属于完全背包问题
// /**
//  * @param {number[]} cost
//  * @param {number} target
//  * @return {number}
//  * @description 给当前结果添加一个数位（i + 1）(1到9)的成本为 cost[i]
//  * 总成本必须恰好等于 target
//  * 添加的数位中没有数字 0 。
//  * @summary 完全背包，选出背包，求最后组成的数的最大值
//  */
const largestNumber = function (cost: number[], target: number): string {
  const dp = Array<string>(target + 1).fill('#')
  dp[0] = ''

  // 选择数位枚举物品
  for (let i = 0; i < 9; i++) {
    // 正序遍历容量
    for (let j = cost[i]; j <= target; j++) {
      const pre = dp[j - cost[i]]
      if (pre === '#') continue
      const next = (i + 1).toString() + pre
      if (compare(next, dp[j])) dp[j] = next
    }
  }

  return dp[target] === '#' ? '0' : dp[target]

  function compare(str1: string, str2: string) {
    return str1.length === str2.length ? str1 > str2 : str1.length > str2.length
  }
}

console.log(largestNumber([4, 3, 2, 5, 6, 7, 2, 5, 5], 9))
// 输出："7772"
// 解释：添加数位 '7' 的成本为 2 ，添加数位 '2' 的成本为 3 。
// 所以 "7772" 的代价为 2*3+ 3*1 = 9 。 "977" 也是满足要求的数字，
// 但 "7772" 是较大的数字。

// 没懂
// 因为最后的答案总要把数字大的放到前面，
// 比如说在7772和7727之间肯定选7772。
// 那么贪心地把先数字大的物品塞到背包里，
// 后面的小数字还能放进去的话就放到数字的最后。
