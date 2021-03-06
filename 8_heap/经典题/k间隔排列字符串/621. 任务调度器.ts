/**
 * @param {string[]} tasks
 * @param {number} n
 * @return {number}
 * @description
 * 苺个任务都可以在 1 个单位时间内执行完
 * 两个 相同种类 的任务之间必须有长度为整数 n 的冷却时间，因此至少有连续 n 个单位时间内 CPU 在执行不同的任务，或者在待命状态。
 * 计算完成所有任务所需要的 最短时间
 *
 * @summary
 * 假设有无数个桶
 * 我们设计桶的大小为 n+1，则相同的任务恰好不能放入同一个桶，最密也只能放入相邻的桶。
 * 一个桶不管是否放满，其占用的时间均为 n+1，这是因为后面桶里的任务需要等待冷却时间。
 * 最后一个桶是个特例，由于其后没有其他任务需等待，所以占用的时间为桶中的任务个数。
 * 总排队时间 = (桶个数 - 1) * (n + 1) + 最后一桶的任务数
 * 如果任务太多桶放不下，则取tasks.length
 */
const leastInterval = (tasks: string[], n: number): number => {
  if (n === 0) return tasks.length

  const counter = new Map<string, number>()
  tasks.forEach(task => counter.set(task, (counter.get(task) || 0) + 1))

  const taskCount = [...counter.values()]
  const maxCount = Math.max.apply(null, taskCount)

  // 数量最大的任务个数
  const maxTaskCount = taskCount.filter(v => v === maxCount).length
  const time = (maxCount - 1) * (n + 1) + maxTaskCount

  // 如果桶放不下，则取tasks.length
  return Math.max(time, tasks.length)
}

console.log(leastInterval(['A', 'A', 'A', 'B', 'B', 'B'], 2))
// A -> B -> (待命) -> A -> B -> (待命) -> A -> B
console.log(leastInterval(['A', 'A', 'A', 'B', 'B', 'B'], 0))
console.log(leastInterval(['A', 'A', 'A', 'B', 'B', 'B'], 0))

export default 1
