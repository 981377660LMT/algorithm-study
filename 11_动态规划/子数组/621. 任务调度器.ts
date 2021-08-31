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
 * 对于重复的任务，我们只能将每个都放入不同的桶中，因此桶的个数就是重复次数最多的任务的个数
 * 一个桶不管是否放满，其占用的时间均为 n+1，这是因为后面桶里的任务需要等待冷却时间。
 * 最后一个桶是个特例，由于其后没有其他任务需等待，所以占用的时间为桶中的任务个数。
 * 总排队时间 = (桶个数 - 1) * (n + 1) + 最后一桶的任务数
 */
const leastInterval = (tasks: string[], n: number): number => {
  const counter = new Map<string, number>()
  tasks.forEach(task => counter.set(task, (counter.get(task) || 0) + 1))
  const vals = [...counter.values()]
  const bucketNum = Math.max.apply(null, vals)
  const lastBucketCount = vals.filter(v => v === bucketNum).length
  const time = (bucketNum - 1) * (n + 1) + lastBucketCount
  return Math.max(tasks.length, time)
}

console.log(leastInterval(['A', 'A', 'A', 'B', 'B', 'B'], 2))
// A -> B -> (待命) -> A -> B -> (待命) -> A -> B
export default 1
