1. 向上标记法 O(n)
2. 倍增在线查询 预处理 nlogn logn
   - bfs 求出 levelMap 和 parentMap
   - 求出 fa[i][j] 从 i 向上走 2^j 步走到的结点
   - 根据 levelMap 求深度
     1. 先将深的点跳到同一层，二进制拼凑出 depth[x]-depth[y]
     2. 两个点同时往上跳 二进制拼凑直到跳到最近公共祖先的下一层
     3. x 或 y 再跳一次即 LCA
3. tarjan 离线求 LCA O(n+m)
