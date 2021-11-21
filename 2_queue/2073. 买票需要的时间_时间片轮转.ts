/**
 * @param {number[]} tickets
 * @param {number} k
 * @return {number}
 * 每个人买票都需要用掉 恰好 1 秒 。一个人 一次只能买一张票
 * 返回位于位置 k（下标从 0 开始）的人完成买票需要的时间（以秒为单位）。
 * k前面的人最多买 tickets[k]张票，k后面的人最多买 tickets[k]-1张票
 */
function timeRequiredToBuy(tickets: number[], k: number): number {
  const need = tickets[k]
  let res = 0
  for (let i = 0; i < tickets.length; i++) {
    if (i <= k) res += Math.min(need, tickets[i])
    else res += Math.min(need - 1, tickets[i])
  }
  return res
}

console.log(timeRequiredToBuy([2, 3, 2], 2))
// 输出：6
// 解释：
// - 第一轮，队伍中的每个人都买到一张票，队伍变为 [1, 2, 1] 。
// - 第二轮，队伍中的每个都又都买到一张票，队伍变为 [0, 1, 0] 。
// 位置 2 的人成功买到 2 张票，用掉 3 + 3 = 6 秒。
