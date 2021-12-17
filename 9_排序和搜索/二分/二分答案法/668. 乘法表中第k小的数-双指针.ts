/**
 * @param {number} m
 * @param {number} n
 * @param {number} k
 * @return {number}
 * 给定高度m 、宽度n 的一张 m * n的乘法表，以及正整数k，你需要返回表中第k 小的数字。
 * 感觉 第k小 第k大的数 简单中等题就是优先队列，难题就是二分
 * nlogn
 * 378. 有序矩阵中第 K 小的元素.ts
 */
const findKthNumber = function (m: number, n: number, k: number): number {
  // 小于等于mid的个数
  const count = (mid: number): number => {
    let count = 0
    // 统计每行符合条件的个数
    for (let i = 1; i <= m; i++) {
      count += Math.min(~~(mid / i), n)
    }

    return count
  }

  let l = 1
  let r = m * n
  while (l <= r) {
    const mid = ~~((l + r) / 2)
    if (count(mid) < k) l = mid + 1
    else r = mid - 1
  }

  return l
}

console.log(findKthNumber(3, 3, 5))
// 输出: 3
// 解释:
// 乘法表:
// 1	2	3
// 2	4	6
// 3	6	9

// 第5小的数字是 3 (1, 2, 2, 3, 3).
export default 1
