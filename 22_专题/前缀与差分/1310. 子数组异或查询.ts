/**
 * @param {number[]} arr 1 <= arr.length <= 3 * 10^4 1 <= arr[i] <= 10^9
 * @param {number[][]} queries 1 <= queries.length <= 3 * 10^4
 * @return {number[]}
 * 请你计算从 Li 到 Ri 的 XOR 值
 * @summary 我们利用了异或的性质 x ^ y ^ x = y
 */
function xorQueries(arr: number[], queries: number[][]): number[] {
  // 前缀异或数组
  const preXOR = Array<number>(arr.length + 1).fill(0)
  for (let i = 1; i < preXOR.length; i++) {
    preXOR[i] = preXOR[i - 1] ^ arr[i - 1]
  }

  return queries.map(([l, r]) => preXOR[l] ^ preXOR[r + 1])
}

console.log(
  xorQueries(
    [1, 3, 4, 8],
    [
      [0, 1],
      [1, 2],
      [0, 3],
      [3, 3],
    ]
  )
)
// 输出：[2,7,14,8]
// 解释：
// 数组中元素的二进制表示形式是：
// 1 = 0001
// 3 = 0011
// 4 = 0100
// 8 = 1000
// 查询的 XOR 值为：
// [0,1] = 1 xor 3 = 2
// [1,2] = 3 xor 4 = 7
// [0,3] = 1 xor 3 xor 4 xor 8 = 14
// [3,3] = 8

export default 1
