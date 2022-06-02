1. 不合题意(障碍、剪枝)
2. 终点
3. 转移

```Python
@lru_cache(None)
def dfs(row: int, col: int) -> List[int]:
   """第一个整数是 「得分」 的最大值，第二个整数是得到最大得分的方案数"""
   if board[row][col] == 'X':  # 不合题意的障碍
       return [-INF, 0]
   if (row, col) == (0, 0):  # 终点
       return [0, 1]

   max_, count = -int(1e20), 0
   for dr, dc in DIR3:   # 转移
       nr, nc = row + dr, col + dc
       if 0 <= nr < ROW and 0 <= nc < COL:
           nextMax, nextCount = dfs(nr, nc)
           maxCand = grid[row][col] + nextMax
           if max_ < maxCand:
               max_, count = maxCand, nextCount
           elif max_ == maxCand:
               count += nextCount
               count %= MOD

   return [max_, count]

```
