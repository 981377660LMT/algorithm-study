二分图应用：婚配问题，宠物收养，工人调配

- 验证二分图

1. dfs 方法:color 数组+从每个未看过的节点开始 dfs+**对每个 visited 的点**验证颜色是否合理
   `785已知邻接表判断二分图`
2. 并查集方法 计算联通分量 count
   `947. 移除最多的同行或同列石头`

二分图(带颜色)模板解法
`886. 可能的二分法`
`1129. 颜色交替的最短路径`

- 无权二部图最大匹配

1. 匈牙利算法
2. 最大流算法(EK/dinic)

- 有权二部图最大权重匹配

1. KM 算法(找权重最大的边组成的子图--------→ 在这个子图上找最大匹配)
2. 最大流算法(EK/dinic)

一棵无向树是二部图
