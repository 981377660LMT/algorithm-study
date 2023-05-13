// 在 manager 数组中，每个员工都有一个直属负责人，其中 manager[i] 是第 i 名员工的直属负责人
// 总负责人，manager[headID] = -1
// 公司总负责人想要向公司所有员工通告一条紧急消息。他将会首先通知他的直属下属们，然后由这些下属通知他们的下属，直到所有的员工都得知这条紧急消息。
// 第 i 名员工需要 informTime[i] 分钟来通知它的所有直属下属（也就是说在 informTime[i] 分钟后，他的所有直属下属都可以开始传播这一消息）。
// 返回通知所有员工这一紧急消息所需要的 分钟数 。

// 即：求出根节点到叶子节点的`最长路径` => dfs
function numOfMinutes(n: number, headID: number, manager: number[], informTime: number[]): number {
  const adjList: [number, number][][] = Array(n)
  for (let i = 0; i < n; i++) adjList[i] = []

  for (let i = 0; i < n; i++) {
    if (i === headID) continue
    const parent = manager[i]
    const weight = informTime[i]
    adjList[parent].push([i, weight])
  }

  let res = 0
  dfs(headID, 0)
  return res + informTime[headID]

  function dfs(cur: number, pathSum: number): void {
    res = Math.max(res, pathSum)
    for (const [next, weight] of adjList[cur]) {
      dfs(next, pathSum + weight)
    }
  }
}

console.log(numOfMinutes(6, 2, [2, 2, -1, 2, 2, 2], [0, 0, 1, 0, 0, 0]))
// 输出：1
// 解释：id = 2 的员工是公司的总负责人，也是其他所有员工的直属负责人，他需要 1 分钟来通知所有员工。
// 上图显示了公司员工的树结构。
