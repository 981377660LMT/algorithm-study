// 从0到n 每个数字表示一个节点
// 如果x 与 y相差一个完全平方数，则连接一条边
// 问题转化为这个有向无权图中从n到0的最短路径
// 树中层序遍历就是广度优先遍历(bfs的queue那一套方法)
// 三步:
// 1.初始化queue和visited
// 2.queue中shift出
// 3.push新的节点
const numSqures = (n: number) => {
  // 节点编号/走了几步
  // 初始值
  const queue: [number, number][] = [[n, 0]]
  const visited = new Set<number>()

  // bfs
  while (queue.length) {
    const [num, steps] = queue.shift()!
    if (num === 0) return steps

    // 重复推入了很多节点，原因：树中一个节点唯一确定一条路径，而图有多种可能
    // 图的遍历需要visited 集合
    for (let i = 0; num - i ** 2 >= 0; i++) {
      const v = num - i ** 2
      !visited.has(v) && queue.push([v, steps + 1])
      visited.add(v)
    }
  }

  throw new Error('No Solution')
}

console.log(numSqures(12))

export {}
