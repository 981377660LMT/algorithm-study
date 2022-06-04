import { useMinCostMaxFlow } from './useMinCostMaxFlow'

function maximumANDSum(nums: number[], numSlots: number): number {
  const n = nums.length
  const m = numSlots
  const [START, END] = [n + m + 2, n + m + 3]
  const mcmf = useMinCostMaxFlow(n + m, START, END)
  for (let i = 0; i < n; i++) {
    for (let j = 0; j < numSlots; j++) {
      mcmf.addEdge(i, j + n, 1, -(nums[i] & (j + 1)))
    }
  }

  for (let i = 0; i < n; i++) mcmf.addEdge(START, i, 1, 0)
  for (let i = 0; i < numSlots; i++) mcmf.addEdge(i + n, END, 2, 0)
  return -mcmf.work()[1]
}

if (require.main === module) {
  console.log(maximumANDSum([1, 2, 3, 4, 5, 6], 3))
}

export {}
