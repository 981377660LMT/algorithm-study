// 设计一个数据结构,有效地找到给定子数组的 出现 threshold 次数或次数以上的元素 。
// 1 <= arr.length <= 2e4
// threshold <= right - left + 1
// 2 * threshold > right - left + 1
// 调用 query 的次数最多为 1e4

// !按照`区间长度`根号分治
// !针对不同的询问区间长度，使用两种不同的方法。
// 记 SQRT = sqrt(2n)
// - 区间长度小于 SQRT ，使用暴力计算
// - 区间长度大于 SQRT ，则绝对众数出现次数 大于 SQRT/2
//   可能的候选人个数不超过 2n/SQRT ，使用前缀和统计频率大于SQRT/2的数的出现次数
