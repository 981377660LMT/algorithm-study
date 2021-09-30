/**
 *
 * @param n 请你构造一个答案列表 answer ，该列表应当包含从 1 到 n 的 n 个不同正整数
 * @param k
 * 相邻两个数差的绝对值有且仅有 k 个不同整数(考虑摆动序列)
 */
function constructArray(n: number, k: number): number[] {
  const res = Array.from({ length: n }, (_, i) => i + 1)

  for (let i = 1, flag = 1; k > 0; i++, flag *= -1, k--) {
    res[i] = res[i - 1] + k * flag
  }
  return res
}

console.log(constructArray(3, 1))
// 1,2,3,4,5 假设4个不同
// 加4减3加2减1
