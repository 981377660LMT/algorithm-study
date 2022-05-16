/**
 * @param {number} n
 * @return {number[]}
 * 给出时间复杂度为O(n*sizeof(integer))的解答非常容易。但你可以在线性时间O(n)内用一趟扫描做到吗
 * @summary
 * 二进制表示中，
 * 奇数一定比前面那个偶数多一个 1，因为多的就是最低位的 1。
 * 偶数中 1 的个数一定和除以 2 之后的那个数一样多
 */
function countBits(n: number): number[] {
  const res = Array<number>(n + 1).fill(0)
  for (let i = 1; i <= n; i++) {
    res[i] = i % 2 === 0 ? res[i / 2] : res[i - 1] + 1
    // 简写:
    res[i] = res[i >> 1] + (i % 2)
  }
  return res
}

console.log(countBits(5))
// 输出: [0,1,1,2,1,2]
export {}
