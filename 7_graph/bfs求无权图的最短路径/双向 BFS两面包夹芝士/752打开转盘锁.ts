/**
 * @param {string[]} deadends
 * @param {string} target
 * @return {number}
 * @description 每次旋转都只能旋转一个拨轮的一位数字
 * @description 给出解锁需要的最小旋转次数(无权图最短路径) 锁的初始数字为 '0000'
 * @description 每次next有八个顶点
 * 双向 BFS 也有局限，因为你必须知道终点在哪里
 * 双向 BFS 还是遵循 BFS 算法框架的，只是不再使用队列，
 * 而是使用 HashSet 方便快速判断两个集合是否有交集。
 *
 * https://labuladong.gitbook.io/algo/mu-lu-ye/bfs-kuang-jia
 */
const openLock = (deadends: string[], target: string): number => {
  const dead = new Set(deadends)
  // 用集合不用队列，可以快速判断元素是否存在
  let queue1 = new Set<string>(['0000'])
  let queue2 = new Set<string>([target])
  const visited = new Set()
  let res = 0

  while (queue1.size && queue2.size) {
    if (queue1.size > queue2.size) {
      ;[queue1, queue2] = [queue2, queue1]
    }
    // 本层搜出来的结果
    const tmp = new Set<string>()
    for (const cur of queue1) {
      if (dead.has(cur)) continue
      if (queue2.has(cur)) return res
      visited.add(cur)

      for (const next of getNextStates(cur)) {
        if (visited.has(next)) continue
        tmp.add(next)
      }
    }

    res++
    queue1 = queue2
    queue2 = tmp
  }
  return -1

  // 八个临边
  function getNextStates(s: string): string[] {
    const ans: string[] = []
    for (let i = 0; i < s.length; i++) {
      ans.push(s.slice(0, i) + ((+s[i] + 1) % 10).toString() + s.slice(i + 1))
      // +9模10表示-1
      ans.push(s.slice(0, i) + ((+s[i] + 9) % 10).toString() + s.slice(i + 1))
    }
    return ans
  }
}

console.log(openLock(['0201', '0101', '0102', '1212', '2002'], '0202'))

export {}
