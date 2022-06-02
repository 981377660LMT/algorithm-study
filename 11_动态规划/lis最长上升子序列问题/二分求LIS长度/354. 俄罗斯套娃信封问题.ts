import { bisectLeft } from '../../../9_排序和搜索/二分/bisect'

/**
 * @link https://leetcode-cn.com/problems/russian-doll-envelopes/solution/ti-mu-zong-jie-zui-chang-shang-sheng-zi-twyhq/
 * @param {number[][]} envelopes
 * @return {number}
 * 当另一个信封的宽度和高度都比这个信封大的时候，这个信封就可以放进另一个信封里
 * @summary 排序之后就是LIS问题
 */
const maxEnvelopes = function (envelopes: number[][]): number {
  if (envelopes.length <= 1) return envelopes.length

  // 保证[6,7]在[6,4]前面 宽相同时不能构成上升序列
  envelopes.sort((a, b) => a[0] - b[0] || b[1] - a[1])
  const LIS: number[] = [envelopes[0][1]]

  for (let i = 1; i < envelopes.length; i++) {
    if (envelopes[i][1] > LIS[LIS.length - 1]) {
      LIS.push(envelopes[i][1])
    } else {
      console.log(LIS, envelopes[i])
      LIS[bisectLeft(LIS, envelopes[i][1])] = envelopes[i][1]
    }
  }

  return LIS.length
}

console.log(
  maxEnvelopes([
    [1, 2],
    [2, 3],
    [3, 4],
    [3, 5],
    [4, 5],
    [5, 5],
    [5, 6],
    [6, 7],
    [7, 8],
  ])
)
// 7
