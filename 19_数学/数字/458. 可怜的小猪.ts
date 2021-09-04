/**
 * @param {number} buckets
 * @param {number} minutesToDie
 * @param {number} minutesToTest
 * @return {number}
 * @description
 * 有 buckets 桶液体，其中 正好 有一桶含有毒药，其余装的都是水
 * 小猪喝完后，必须有 minutesToDie 分钟的冷却时间。
 * 你只有 minutesToTest 分钟时间来确定哪桶液体是有毒的
 * 返回在规定时间内判断哪个桶有毒所需的 最小 猪数
 * @summary 死亡状态
 * 只能喝一次 => 二进制 (死、活)
 * 只能喝两次 => 三进制 (死0次 死1次 死2次)
 * ...
 * @link https://leetcode-cn.com/problems/poor-pigs/comments/107570
 */
const poorPigs = (buckets: number, minutesToDie: number, minutesToTest: number): number => {
  if (buckets === 1) return 0
  const radix = ~~(minutesToTest / minutesToDie) + 1
  // 桶的编号从0开始
  return (buckets - 1).toString(radix).length
}

console.log(poorPigs(1000, 15, 60))
// 输出：5
export default 1

// 假设是1000个桶我可以用二进制标记 从0000000001 到 1111101000。
// 对应一共10位（从右到左标记分别1到10位）。
// 10只猪每只🐖喝对应位置上是1混合水。
// 最后看死了那几只猪则对应位置则为1. 假设是60分钟。
// 因为每只猪在规定时间可以喝4次水。
// 因此是5进制表示看有多少位则需要几只猪。1000是有5位因此是5只猪。
