import { memo } from '../../5_map/memo'

function minimumWhiteTiles(floor: string, numCarpets: number, carpetLen: number): number {
  const dfs = memo((index: number, remain: number): number => {
    if (index >= n || remain === 0) return 0
    let res = 0
    res = Math.max(res, dfs(index + 1, remain))
    res = Math.max(
      res,
      dfs(index + carpetLen, remain - 1) + preSum[Math.min(n, index + carpetLen)] - preSum[index]
    )
    return res
  })

  const n = floor.length
  const nums = floor.split('').map(Number)
  const preSum = Array<number>(n + 1).fill(0)
  for (let i = 0; i < n; i++) preSum[i + 1] = preSum[i] + nums[i]
  const ones = nums.filter(x => x === 1).length

  const res = dfs(0, numCarpets)
  dfs.cacheClear()
  return ones - res
}
