- 图的表示
- 图的遍历
- dfs:连通性/找路径/二分图检测/环的检测/floodfill
- bfs:无权图最短路径
- 使用图论建模

地铁路线规划
互联网连接 Google PageRank
爬虫 图(页面)的遍历
论文引用
社交网络
匹配问题
编译原理 babel
状态机

2. 图的基本表示(建立图的顶点到实际对象的双向映射)
3. dfs 检测环:递归的函数利用**函数返回值表示不同的意义**，例如返回 false 代表无环，由此更好地递归 不是二分图/有环则马上 **return false/true 中止判断**
4. 图论问题建模
   生成迷宫:**随机队列** push 点 经过的点染色
5. bfs 模板

```TS
const solution = () => {
  // pre 存路径
  const pre = new Map<number, number>()
  const queue: number[] = [0]
  const visited = new Set([0])
  let end = -1

  const bfs = () => {
    while (queue.length) {
      const cur = queue.shift()!
      // 1. if(结束条件) return  end=...

      // 2. ... 求下一个状态

      for (const next of 下一个状态) {
        if (!visited.has(next)) {
          queue.push(next)
          pre.set(next, cur)
          visited.add(next)
        }
      }
    }
  }
  bfs()

  // 3.从路径倒推
  let p = end
  const res: number[] = [end]
  while (pre.get(p)) {
    p = pre.get(p)!
    res.push(p)
  }
  res.push(0)

  return res.reverse()
}

```

6. 寻找桥的算法使用 DFS 就可以解决，其实就是 Tarjan 算法
   需要对每一条边判断
   如何判断 v-w 不是桥? 看通过 w，能否从另外一条路回到 v 或 v 之前的顶点
   需要记录的两个信息 order(dfs 的顺序) 和 low(记录这个顶点所属的组的序号(最小的 order)，不断遍历比较然后回退比较两个节点，取为每个连通区域的最小值 order)
   如果 low[w]大于 order[v]则说明是 v-w 桥 low 相等则同一组

7. 注意 dfs 传参 是使用全局的变量还是每个 dfs 一个变量
   例如 count 这种变量应该传参而不是全局 每个 dfs 都有一个 count

8. 寻找欧拉回路
   fleury 算法(贪心算法 有多条边的时候不走桥 对每一个临边判断是不是桥 不能预处理) 复杂度 O((V+E)^2)
   **hierholzer 算法**
   使用栈保存历史记录:从一个点出发，如果有 next,走到 next 并删除原来走过的边，入栈；否则记入 res 中,栈弹出并回退
   curPath 栈寻找环 入栈 走到不能再走了回溯出栈 继续寻找环入栈 出栈
   每个边走一次 回退一次 O(E)

```JS
  get eulerLoop(): number[] {
    if (!this.hasEulerLoop) return []
    const res: number[] = []
    const clonedAdjMap = this.adjMap.cloneAdjMap()
    let cur = 0
    const stack: number[] = [cur]

    while (stack.length) {
      if (clonedAdjMap.degree(cur) !== 0) {
        stack.push(cur)
        const next = clonedAdjMap.adj(cur).shift()!
        clonedAdjMap.removeEdge(cur, next)
        cur = next
      } else {
        // 回退
        res.push(cur)
        cur = stack.pop()!
      }
    }

    return res
  }
```

9. 无向带权图最小生成树问题

无权图的邻接表:

```TS
number[][]
[[0,1,2],[3],[1,2],[0,1]]
```

带权图的邻接表:

```TS
Map<number,number>[]
[
  Map(3) { 1 => 2, 3 => 7, 5 => 2 },
  Map(5) { 0 => 2, 2 => 1, 3 => 4, 4 => 3, 5 => 5 },
  Map(3) { 1 => 1, 4 => 4, 5 => 4 },
  Map(4) { 0 => 7, 1 => 4, 4 => 1, 6 => 5 },
  Map(4) { 1 => 3, 2 => 4, 3 => 1, 6 => 7 },
  Map(3) { 0 => 2, 1 => 5, 2 => 4 },
  Map(2) { 3 => 5, 4 => 7 }
]
```

10. 动态判断是否构成环:并查集
    chapter11 无向带权图\带权图\最小生成树.ts

| Kruskal | O(ElogE) | 借助并查集           |
| ------- | -------- | -------------------- |
| Prim    | O(ElogE) | 借助优先队列(最小堆) |

11. 最短路径算法
1. Dijkstra 找最小值=>确定解=>更新
