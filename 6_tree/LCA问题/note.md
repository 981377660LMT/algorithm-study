## 查询两个树结点的 LCA

1. 离线查询

   - **tarjan** O(n+q) dfs+并查集

2. 在线查询

   - **倍增** O(nlogn)预处理 O(logn)查询，两个先跳到一样的高度再同时往上跳直到相等
   - **线性** O(n)

     1. 上跳记录访问的路径，第一个交点即为 LCA
     2. 相交链表找交点(空间复杂度 O(1))
     3. 预处理 depth 和 parent 上跳，两个先跳到一样的高度再同时往上跳直到相等
     4. 从根节点自顶向下递归查询

   - **树链剖分**，O(n)预处理 O(logn)查询，沿着最多 logn 段链上跳
   - **dfs 序(欧拉序) + st 表**， O(nlogn)预处理 O(1)查询，这两点之间的区间中，深度最小点就是 LCA。这可以用 RMQ 解决
     > https://www.cnblogs.com/pealicx/p/6859901.html

3. 支持在线加点删点
   - **link-cut tree** O(nlogn)预处理 O(logn)查询 ，树链剖分 + splay 维护链的信息
   - **toptree** O(nlogn)预处理 O(logn)查询
