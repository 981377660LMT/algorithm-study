import { PriorityQueue } from '../todo优先级队列'

type Start = number
type Cost = number
type Index = number
/**
 * @param {number[][]} tasks  [enqueueTimei, processingTimei] 意味着第 i​​​​​​​​​​ 项任务将会于 enqueueTimei 时进入任务队列，需要 processingTimei 的时长完成执行。
 * @return {number[]}  返回 CPU 处理任务的顺序。
 * @summary
 * 我们可以直接模拟即可。
 */
const getOrder = function (tasks: number[][]): number[] {
  const res: number[] = []

  // 进入时间，下标，耗费时间
  const sortedTasks = tasks
    .map<[Start, Cost, Index]>(([start, cost], i) => [start, cost, i])
    .sort((a, b) => a[0] - b[0] || a[1] - b[1] || a[2] - b[2])

  const pq = new PriorityQueue<[Cost, Index]>((a, b) => a[0] - b[0] || a[1] - b[1])
  let timeStamp = 0
  let nextTaskIndex = 0

  // 一直循环
  while (res.length < tasks.length) {
    // 1.放入所有合理的任务
    while (nextTaskIndex < sortedTasks.length && sortedTasks[nextTaskIndex][0] <= timeStamp) {
      pq.push([sortedTasks[nextTaskIndex][1], sortedTasks[nextTaskIndex][2]])
      nextTaskIndex++
    }

    // 2. 执行一个任务
    if (pq.length) {
      const [cost, index] = pq.shift()!
      res.push(index)
      timeStamp += cost
    } else {
      // 3.队列为空 则时间跳到下一个任务开始
      timeStamp = sortedTasks[nextTaskIndex][0]
    }
  }

  return res
}

console.log(
  getOrder([
    [1, 2],
    [2, 4],
    [3, 2],
    [4, 1],
  ])
)
