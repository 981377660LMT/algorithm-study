求第 n 次变换的第 k 个数
**类似线段树**
`dfs(depth,target)`
求出 mid 后，比较 target 和 `mid` 的关系
target <= mid 向左递归
target > mid 向右递归
