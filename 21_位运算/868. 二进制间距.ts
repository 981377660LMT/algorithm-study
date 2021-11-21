function binaryGap(n: number): number {
  let preIndex = Infinity
  let res = 0

  for (let i = 0; i < 32; i++) {
    if (((n >> i) & 1) === 1) {
      if (preIndex !== Infinity) res = Math.max(res, i - preIndex)
      preIndex = i
    }
  }

  return res
}
// 输入：n = 22
// 输出：2
// 解释：
// 22 的二进制是 "10110" 。
// 在 22 的二进制表示中，有三个 1，组成两对相邻的 1 。
// 第一对相邻的 1 中，两个 1 之间的距离为 2 。
// 第二对相邻的 1 中，两个 1 之间的距离为 1 。
// 答案取两个距离之中最大的，也就是 2 。
