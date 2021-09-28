type Min = number
type Max = number
/**
 * @param {number[]} nums
 * @return {number}
 * 给定一个无序的数组，找出数组在排序之后，相邻元素之间最大的差值。
 * 请尝试在线性时间复杂度和空间复杂度的条件下解决此问题。(暗示桶排序)
 * @summary
 * 设长度为 N 的数组中最大最小值为 max,min，
 * 则不难发现相邻数字的最大间距不会小于(max-min)/(N-1) (反证)

 */
var maximumGap = function (nums: number[]): number {
  if (nums.length < 2) return 0
  const min = Math.min.apply(null, nums)
  const max = Math.max.apply(null, nums)
  const size = Math.max(1, ~~((max - min) / (nums.length - 1)))
  const bucketNum = ~~((max - min) / size) + 1

  const bucket = Array.from<unknown, [Min, Max]>({ length: bucketNum }, () => [-1, -1])
  for (const num of nums) {
    const index = ~~((num - min) / size)
    if (bucket[index][0] === -1) {
      bucket[index][0] = bucket[index][1] = num
    } else {
      bucket[index][0] = Math.min(bucket[index][0], num)
      bucket[index][1] = Math.max(bucket[index][1], num)
    }
  }

  // 每个桶跟相邻桶比较即可 当前桶min-之前桶max
  let res = 0
  let preBucket = -1
  // console.log(bucket, bucketNum)
  for (let i = 0; i < bucketNum; i++) {
    if (bucket[i][0] === -1) continue
    if (preBucket !== -1) res = Math.max(res, bucket[i][0] - bucket[preBucket][1])
    preBucket = i
  }
  return res
}

console.log(maximumGap([3, 6, 9, 1]))
console.log(maximumGap([1, 10]))
