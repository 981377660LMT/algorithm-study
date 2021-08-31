// 利用计数排序的思想O(n)
// 全部对号入座 计算连续0的最大个数
const maxDiff = (nums: number[]) => {
  const max = Math.max.apply(null, nums)
  const min = Math.min.apply(null, nums)
  const bucket = Array<number>(max - min + 1).fill(0)
  nums.forEach(num => {
    bucket[num - min] += 1
  })
  const tmp = bucket.join('')
  const match = (tmp.match(/(0)\1*/g) || []).map(v => v.length)
  return Math.max.apply(null, match) + 1
}

console.log(maxDiff([2, 6, 3, 4, 5, 10, 9]))
