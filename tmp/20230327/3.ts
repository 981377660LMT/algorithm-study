import { WaveletMatrixSum } from '../../24_高级数据结构/waveletmatrix/WaveletMatrixSum'

function minOperations(nums: number[], queries: number[]): number[] {
  const wm = new WaveletMatrixSum(new Uint32Array(nums))
  const n = nums.length
  const res: number[] = []
  const sum = nums.reduce((a, b) => a + b, 0)
  queries.forEach(target => {
    const [leftCount, leftSum] = wm.countPrefix(0, n, target)
    const rightSum = sum - leftSum
    const rightCount = n - leftCount
    res.push(leftCount * target - leftSum + rightSum - rightCount * target)
  })
  return res
}
