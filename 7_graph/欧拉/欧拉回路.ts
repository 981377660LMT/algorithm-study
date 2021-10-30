/**
   * @description Hierholzer算法(插入回路法)
   * （1）选择任一顶点为起点，入栈curPath，深度搜索访问顶点，将经过的边都删除，经过的顶点入栈curPath。
    （2）如果当前顶点没有相邻边，则将该顶点从curPath出栈到res。
    （3）res中的顶点逆序，就是从起点出发的欧拉回路。(当然顺序也是)
   */
function eulerLoop(adjList: number[][]): number[] {
  const res: number[] = []
  let cur = 0
  const stack: number[] = [cur]

  while (stack.length) {
    if (adjList[cur].length !== 0) {
      stack.push(cur)
      const next = adjList[cur].pop()!
      cur = next
    } else {
      // 回退
      res.push(cur)
      cur = stack.pop()!
    }
  }

  // 有向图需要反转
  return res.reverse()
}

console.log(
  eulerLoop([
    [1, 2],
    [0, 3],
    [1, 2],
    [3, 0],
  ])
)
