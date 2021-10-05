/**
 * @param {number} n
 * @return {number}
 * 假定每分钟选择以下两种策略之一:

   使用当前带宽下载插件
   将带宽加倍（下载插件数量随之加倍）
   请返回小扣完成下载 n 个插件最少需要多少分钟。
   
 * 贪心，优先扩大带宽，最后的加一是下载的那一次操作
 */
var leastMinutes = function (n) {
  return Math.ceil(Math.log(n) / Math.log(2)) + 1
}

console.log(leastMinutes(4))
// 输入：n = 4

// 输出：3

// 解释：
// 最少需要 3 分钟可完成 4 个插件的下载，以下是其中一种方案:
// 第一分钟带宽加倍，带宽可每分钟下载 2 个插件;
// 第二分钟下载 2 个插件;
// 第三分钟下载 2 个插件。
