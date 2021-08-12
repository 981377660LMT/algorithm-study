```Python
seen = set()
def dfs(i, j):
  if i 越界 or j 越界: return
  if (i, j) in seen: return
  temp = board[i][j]
  # 标记为访问过
  seen.add((i, j))
  # 上
  dfs(i + 1, j)
  # 下
  dfs(i - 1, j)
  # 右
  dfs(i, j + 1)
  # 左
  dfs(i, j - 1)
  # 撤销标记
  seen.remove((i, j))
# 单点搜索
dfs(0, 0)
# 多点搜索
for i in range(M):
   for j in range(N):
      dfs(i, j)
```

原地标记

```Python
def dfs(i, j):
  if i 越界 or j 越界: return
  if board[i][j] == -1: return
  temp = board[i][j]
  # 标记为访问过
  board[i][j] = -1
  # 上
  dfs(i + 1, j)
  # 下
  dfs(i - 1, j)
  # 右
  dfs(i, j + 1)
  # 左
  dfs(i, j - 1)
  # 撤销标记
  board[i][j] = temp
# 单点搜索
dfs(0, 0)
# 多点搜索
for i in range(M):
   for j in range(N):
      dfs(i, j)
```
