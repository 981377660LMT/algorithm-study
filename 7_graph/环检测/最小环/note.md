求图中的最小环/最大环的一些方法

1. 无权图最小环:bfs
2. 带权图最小环:
   - floyd O(n^3)
   - dijkstra O(E^2logE)
     Dijkstra/spfa:枚举所有边，删除这条边之后以这条边的端点为起点终点跑一次 Dijkstra
3. 最大环:寻找所有环的分组

   - 拓扑排序+dfs 找到环分组(基环树)
   - Tarjan 缩点成树/拓扑图
     `注意 Tarjan 缩点为了变为拓扑图/树，可以判断每个连通分量内点的个数，但这不是必要的`
     注意并查集只能维护组的连通性/大小，**无法维护环的信息**

---

update:

4. 有向图无向图统一的最小环求法 - 最短路径树
   https://yukicoder.me/problems/no/1320/editorial
   `O(V*(V+E)*logV)`
