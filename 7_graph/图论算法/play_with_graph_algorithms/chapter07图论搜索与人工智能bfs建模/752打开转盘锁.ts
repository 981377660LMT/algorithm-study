/**
 * @param {string[]} deadends
 * @param {string} target
 * @return {number}
 * @description 每次旋转都只能旋转一个拨轮的一位数字
 * @description 给出解锁需要的最小旋转次数(无权图最短路径) 锁的初始数字为 '0000'
 * @description 每次next有八个顶点
 */
const openLock = (deadends: string[], target: string): number => {
  let res = -1
  const dead = new Set(deadends)
  const queue: [string, number][] = [['0000', 0]]
  const visited = new Set(['0000'])
  // 八个临边
  const getNextStates = (s = '0000'): string[] => {
    const ans: string[] = []
    for (let i = 0; i < s.length; i++) {
      ans.push(s.slice(0, i) + ((+s[i] + 1) % 10).toString() + s.slice(i + 1))
      // +9模10表示-1
      ans.push(s.slice(0, i) + ((+s[i] + 9) % 10).toString() + s.slice(i + 1))
    }
    return ans
  }

  const bfs = () => {
    while (queue.length) {
      const [cur, step] = queue.shift()!
      if (cur === target) return (res = step)
      if (dead.has(cur)) continue

      for (const next of getNextStates(cur)) {
        if (!visited.has(next)) {
          visited.add(next)
          queue.push([next, step + 1])
        }
      }
    }
  }
  bfs()

  return res
}

console.log(openLock(['0201', '0101', '0102', '1212', '2002'], '0202'))

export {}
