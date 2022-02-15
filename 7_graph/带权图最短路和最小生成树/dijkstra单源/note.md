Dijkstra 算法的签名：
输入一幅图和一个起点 start，计算 start 到其他节点的最短距离

```JAVA
int[] dijkstra(int start, List<Integer>[] graph);
```

标准的 Dijkstra 算法会把从起点 start 到所有其他节点的最短路径都算出来。

1. 没有 visited 集合记录已访问的节点，所以一个节点会被访问多次，会被多次加入队列，那会不会导致队列永远不为空，造成死循环？
   **这个算法不用 visited 集合也不会死循环**
   循环结束的条件是队列为空，那么你就要注意看什么时候往队列里放元素（调用 offer）方法，再注意看什么时候从队列往外拿元素（调用 poll 方法）。
   while 循环每执行一次，都会往外拿一个元素，但想往队列里放元素，可就有很多限制了，必须满足下面这个条件：
   ```JAVA
   // 看看从 curNode 达到 nextNode 的距离是否会更短
   if (distTo[nextNodeID] > distToNextNode) {
    // 更新 dp table
    distTo[nextNodeID] = distToNextNode;
    pq.offer(new State(nextNodeID, distToNextNode));
   }
   ```
2. 为什么用优先级队列 PriorityQueue 而不是 LinkedList 实现的普通队列？为什么要按照 distFromStart 的值来排序？
   如果你非要用普通队列，其实也没问题的，你可以直接把 PriorityQueue 改成 LinkedList，也能得到正确答案，但是效率会低很多。
   Dijkstra 算法使用优先级队列，主要是为了效率上的优化，类似一种贪心算法的思路。
   Bellman-Ford 算法，这个算法是一种更通用的最短路径算法，因为它可以处理带有负权重边的图，Bellman-Ford 算法逻辑和 Dijkstra 算法非常类似，用到的就是普通队列
3. 如果我只想计算起点 start 到某一个终点 end 的最短路径，是否可以修改算法，提升一些效率？
   shift 处的 nodeid 为 end 则 return
   否则 到达不了

**精髓**
`743. 网络延迟时间` 标准模板
`1631. 最小体力消耗路径Dijkstra一般速度`

1. pq 的建立
2. dist 数组(dp 的作用)
3. bfs 的 next 步骤 中如何向 pq 中加入新的点(更新到原点的距离)
   **模板**

```JS
  // 0.建图(也可以不建，只是方便获取next)
  const adjList = Array.from<number, [number, number][]>({ length: n  }, () => [])
  times.forEach(([u, v, w]) => adjList[u].push([v, w]))

  // 1.dist数组
  const dist = Array<number>(n).fill(Infinity)
  dist[start] = 0

  const visited = new Set<number>()

  // 2.pq优先队列
  const compareFunction = (a: Edge, b: Edge) => a[1] - b[1]
  const priorityQueue = new PriorityQueue<Edge>(compareFunction)
  priorityQueue.push([start, 0])

  while (priorityQueue.length) {
    // 3.每次都从离原点最近的没更新过的点开始更新(性能瓶颈：可使用优先队列优化成ElogE)
    const [cur, maxWeight] = priorityQueue.shift()!
    if (cur === end) return maxWeight

    if (visited.has(cur)) continue
    visited.add(cur)

    // 4.利用cur点来更新其相邻节点next与原点的距离
    for (const [next, weight] of adjList[cur]) {
      if (dist[cur] + weight < dist[next]) {
        dist[next] = dist[cur] + weight
        priorityQueue.push([next, dist[next]])
      }
    }
  }

  return Infinity

```

4. 如果我们要打印具体路径呢？
   其实很简单，我们只需要用一个` pre[]数组存储每个点的父节点即可`
   单源最短路的起点是固定的，所以每条路有且仅有一个祖先节点，一步步溯源上去的路径是唯一的
   每当更新一个点的 dist 时，顺便更新一下它的 pre。

```JS
// 4.利用cur点来更新其相邻节点next与原点的距离
    for (const [next, weight] of adjList[cur]) {
      if (dist[cur] + weight < dist[next]) {
        dist[next] = dist[cur] + weight
        pre[next]=cur  // 这里
        priorityQueue.push([next, dist[next]])
      }
    }
```
