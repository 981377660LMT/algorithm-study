有时用 dfs 会超时 但是仍然需要后序遍历(用子节点更新父结点)

两步：

1. 用 bfs 建立 parent 关系和 order 关系
   `如果出队列(栈)用 pop ，(叶子节点的)order 就是 dfs 序`
   `如果出队列用 popleft ，order 就是 拓扑序`
   [dfs 序](E%20-%20Ranges%20on%20Tree.py)
2. 沿着 order 倒着更新父结点
