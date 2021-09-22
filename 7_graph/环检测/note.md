无向图环检测:

- dfs:走回了之前走过的非 pre 的节点
- 并查集:不断 union union 前检查 isConnected

有向图环检测

- dfs:检测有向图的环,需要记录路径(path 集合)。回溯时删除
- 拓扑排序:ok
