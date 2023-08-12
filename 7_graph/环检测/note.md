**找环请认准拓扑排序**

无向图环检测:

- 拓扑排序:度数为 `1` 的点加入队列，每次碰到一个结点就标记为不在环上
- dfs:走回了之前走过的非 pre 的节点 `dfs(cur,pre):boolean`
- 并查集:不断 union union 前检查 isConnected

有向图环检测

- 拓扑排序:度数为 `0` 的点加入队列，每次碰到一个结点就标记为不在环上
- dfs:如果 onPath 则有环 如果 visited 则无环 `dfs(cur):boolean`

---

有向图在线环检测 O(n^2logn)
Incremental Cycle Detection in Directed Graphs
https://github.com/blutorange/js-incremental-cycle-detect/blob/master/src/GenericGraphAdapter.ts
