// # https://leetcode-cn.com/problems/t3fKg1/solution/10xing-jie-jue-zhan-dou-by-foxtail-ke2e/
// # LCP 32. 批量处理任务
// # 2 <= tasks.length <= 10^5
// # 0 <= tasks[i][0] <= tasks[i][1] <= 10^9
// # 某实验室计算机待处理任务以 [start,end,period] 格式记于二维数组 tasks，
// # 表示完成该任务的时间范围为起始时间 start 至结束时间 end 之间，需要计算机投入 period 的时长
// # !处于开机状态的计算机可同时处理任意多个任务，请返回电脑最少开机多久，可处理完所有任务。

import { DualShortestPath } from '../差分约束'

function processTasks(tasks: number[][]): number {
  const allNums = new Set<number>()
  tasks.forEach(([s, e]) => {
    allNums.add(s - 1).add(e)
  })
  const nums = [...allNums].sort((a, b) => a - b)
  const n = nums.length
  const mp = new Map<number, number>()
  nums.forEach((num, i) => mp.set(num, i))

  const D = new DualShortestPath(n + 10, true)
  tasks.forEach(([s, e, p]) => {
    const u = mp.get(s - 1)!
    const v = mp.get(e)!
    D.addEdge(u, v, -p)
  })
  for (let i = 1; i < n; i++) {
    D.addEdge(i - 1, i, 0)
    D.addEdge(i, i - 1, nums[i] - nums[i - 1])
  }

  const [res, ok] = D.run()
  if (!ok) return -1
  return res[n - 1]
}

if (require.main === module) {
  console.log(
    processTasks([
      [1, 3, 2],
      [2, 5, 3],
      [5, 6, 2]
    ])
  )
}
