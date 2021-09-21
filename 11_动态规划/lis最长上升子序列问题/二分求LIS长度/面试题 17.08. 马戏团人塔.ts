import { bisectLeft } from '../../../9_排序和搜索/二分api/7_二分搜索寻找最左插入位置'

/**
 * @param {number[]} height
 * @param {number[]} weight
 * @return {number}
 * 升序排序身高， 若身高相同，体重按降序排序。
 * 在上面的人要比下面的人矮一点且轻一点。已知马戏团每个人的身高和体重，请编写代码计算叠罗汉最多能叠几个人。
 */
const bestSeqAtIndex = function (height: number[], weight: number[]): number {
  if (weight.length <= 1) return weight.length
  const envelopes = Array.from<number, [number, number]>({ length: height.length }, (_, i) => [
    height[i],
    weight[i],
  ])

  if (envelopes.length <= 1) return envelopes.length

  // 保证[6,7]在[6,4]前面 保证后面的比较逻辑成立
  envelopes.sort((a, b) => a[0] - b[0] || b[1] - a[1])
  const LIS: number[] = [envelopes[0][1]]

  for (let i = 1; i < envelopes.length; i++) {
    if (envelopes[i][1] > LIS[LIS.length - 1]) {
      LIS.push(envelopes[i][1])
    } else {
      LIS[bisectLeft(LIS, envelopes[i][1])] = envelopes[i][1]
    }
  }

  return LIS.length
}

console.log(
  bestSeqAtIndex(
    [2868, 5485, 1356, 1306, 6017, 8941, 7535, 4941, 6331, 6181],
    [5042, 3995, 7985, 1651, 5991, 7036, 9391, 428, 7561, 8594]
  )
)
// 输出：5
