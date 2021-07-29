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
const solution = (start:number) => {
  // pre 存路径
  const pre = new number[] = Array(V).fill(-1)
  const queue: number[] = []
  // 可以用pre代替
  // const visited = new Set([start])
  pre[start]=start
  queue.push(start)

  const bfs = (start:number) => {
    while (queue.length) {
      const cur = queue.shift()!
      // 1. if(结束条件) return  end=...

      // 2. ... 求下一个状态

      for (const next of 下一个状态) {
        if (!visited.has(next)) {
          queue.push(next)
          pre[next]=cur
          // visited.add(next)
        }
      }
    }
  }
  bfs(start)

  // 3.从路径倒推
  let p = end
  const res: number[] = []
  while (pre[p]!==start) {
    res.push(p)
    p = pre[p]
  }
  res.push(start)

  return res.reverse()
}

```

从 pre 数组求路径

```TS
  private getAugPath(pre: number[], start: number, end: number) {
    const res: number[] = []
    let p = end

    while (pre[p] !== start) {
      res.push(p)
      p = pre[p]
    }
    res.push(start)

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
    Dijkstra 找最小值=>确定解=>更新

12. 有向图环检测:程序模块引用/任务调度/学习计划依赖
    无向图检测有环:重复到达即有环
    **有向图检测有环**: 1.重复到达不一定有环，要保证重复到达的点在现在的路径上，即 dfs 回溯时去除 path 中的点 2.拓扑排序

    有向无环图:DAG

两个有向图算法:
**拓扑排序算法**:入度+队列，顺便做环检测 O(V+E)
学了 1 号课程才能学 2 号课程
对入度数组排序，**入度为 0 的点开始遍历**,这个点的 next 入度减一
如果入度降将为 0， 则将节点入队
排序结果不唯一
**如果拓扑排序无解则说明有向图有环**

```TS
  topoSort() {
    let hasCycle = false
    const res: number[] = []
    const queue: number[] = this.dfs.adjMap.outDegrees.filter(inDegree => inDegree === 0)
    console.log(this.dfs.adjMap)
    while (queue.length) {
      const v = queue.shift()!
      res.push(v)
      this.dfs.adjMap.adj(v).forEach(w => {
        this.dfs.adjMap.outDegrees[w]--
        this.dfs.adjMap.outDegrees[w] === 0 && queue.push(w)
      })
    }

    if (res.length < this.dfs.adjMap.V) {
      hasCycle = true
      res.splice(0)
    }

    return { res, hasCycle }
  }
```

**求解有向图强连通分量(Strong Critial Connection)**
将所有强连通分量看做一个点，得到的有向图一定是 DAG
**kosaraju 算法**
反图的 DFS 后序的逆序做 CC

https://www.bilibili.com/video/BV1Pv41157xh/?spm_id_from=333.788.recommend_more_video.0

13. **最大流算法**:
    FF 算法:加反向路径给一个反悔的机会 FF 算法很慢 最坏循环次数是网络流的大小边数 即 O(`f*E`)
    EK 算法 O(VE^2):FF 算法的特例 第一步寻找最短路径(无权图最短路径) 在残量图找增广路径(用 BFS)
    最大流算法建模解决分配问题:棒球比赛
    容量:

    1. 对于每种状态有多少场胜利需要分配(离开源点的容量),
    2. 胜利流向哪个队(中间交叉的容量),
    3. 每个队还能允许多少场胜利(流向汇点的容量)

    看最大流是多少，**能不能接纳从源点的流量总合**

14. **二分图(BipartiteGraph)**匹配问题两种方法
    最大匹配/完全匹配
    **最大流算法**(MaxFlow):新建虚拟的源点和汇点，无向图转换为有向图,所有边容量为 1,流入汇点的每一个流量即为一个匹配
    **匈牙利算法**(Hangurian)：从 左侧非匹配点出发 从右向左走永远走匹配边 终止于另外一个非匹配点:增广路径
    则最大匹配数加一
    匈牙利算法 BFS/DFS **关键** O(VE)

说明:使用 **DFS**这个类 读取图
