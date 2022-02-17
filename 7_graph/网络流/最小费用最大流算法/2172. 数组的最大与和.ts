import { mincostMaxflow } from './MincostMaxflow'

function maximumANDSum(nums: number[], numSlots: number): number {
  const n = nums.length
  const mcmf = mincostMaxflow(n + numSlots + 2)
  const from = n + numSlots
  const to = from + 1
  for (let i = 0; i < n; i++) {
    mcmf.addEdge(from, i, 0, 1)
    for (let j = 0; j < numSlots; j++) {
      mcmf.addEdge(i, j + n, -(nums[i] & (j + 1)), 1)
    }
  }

  for (let i = 0; i < numSlots; i++) mcmf.addEdge(i + n, to, 0, 2)
  return -mcmf.mincostFlow(from, to, 220)
}

console.log(maximumANDSum([1, 2, 3, 4, 5, 6], 3))
