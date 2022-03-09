/**
   * @description Hierholzer算法(插入回路法)
   * （1）选择任一顶点为起点，入栈curPath，深度搜索访问顶点，将经过的边都删除，经过的顶点入栈curPath。
    （2）如果当前顶点没有相邻边，则将该顶点从curPath出栈到res。
    （3）res中的顶点逆序，就是从起点出发的欧拉回路。(当然顺序也是)
    @summary 保证存在欧拉路径
   */
function getEulerPath(
  adjMap: Map<number, Set<number>>,
  start: number,
  isDirected: boolean
): number[] {
  let cur = start
  const stack: number[] = [start]
  const res: number[] = []

  while (stack.length > 0) {
    if (adjMap.has(cur) && adjMap.get(cur)!.size > 0) {
      stack.push(cur)
      const next = adjMap.get(cur)!.keys().next().value!
      if (!isDirected) adjMap.get(next)!.delete(cur)
      cur = next
    } else {
      res.push(cur)
      cur = stack.pop()!
    }
  }

  // 有向图需要反转
  return res.reverse()
}

export {}
