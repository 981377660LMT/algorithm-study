无向图环检测:

- dfs:走回了之前走过的非 pre 的节点 `dfs(cur,pre):boolean`
- 并查集:不断 union union 前检查 isConnected

有向图环检测

- dfs:如果 onPath 则有环 如果 visited 则无环 `dfs(cur):boolean`
- 拓扑排序:ok
