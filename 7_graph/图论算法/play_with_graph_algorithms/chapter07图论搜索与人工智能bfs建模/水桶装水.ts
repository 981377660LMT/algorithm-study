/**
 * @description 有两个水桶,一个装5升，一个装3升。怎么利用这两个水桶，得到4升水? (推广到x,y,x)
 * @description 状态转移有哪些情况 6种
 * @description x,y 即为10*x+y 状态压缩
 * @description 狼🐏菜农夫的问题也可这种思路
 */
const solution = () => {
  const pre = new Map<number, number>()
  const queue: number[] = [0]
  const visited = new Set([0])
  let end = -1

  const bfs = () => {
    while (queue.length) {
      const cur = queue.shift()!
      const a = Math.floor(cur / 10)
      const b = cur % 10

      const aToB = Math.min(a, 3 - b)
      const bToA = Math.min(b, 5 - a)
      const nexts: number[] = [
        5 * 10 + b,
        a * 10 + 3,
        10 * a + 0,
        0 * 10 + b,
        (a - aToB) * 10 + (b + aToB),
        (a + bToA) * 10 + (b - bToA),
      ]

      for (const next of nexts) {
        if (!visited.has(next)) {
          queue.push(next)
          pre.set(next, cur)
          visited.add(next)
          if (Math.floor(next / 10) === 4 || next % 10 === 4) return (end = next)
        }
      }
    }
  }
  bfs()

  let p = end
  const res: number[] = [end]
  while (pre.get(p)) {
    p = pre.get(p)!
    res.push(p)
  }
  res.push(0)

  return res.reverse()
}

console.log(solution())

export {}
