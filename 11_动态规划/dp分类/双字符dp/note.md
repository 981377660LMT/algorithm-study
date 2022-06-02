dp[i][j]:=第一串考虑[0..i]，第二串考虑[0..j]时，原问题的解
其中一种最常见的状态转移形式：推导 dp[i][j] 时，dp[i][j] 仅与 dp[i-1][j], dp[i][j-1], dp[i-1][j-1]

```Python
@lru_cache(None)
  def dfs(i: int, j: int) -> bool:
      if i == len(s1) and j == len(s2):
          return True
      res = False
      if i < len(s1) and s1[i] == s3[i + j]:
          res = res or dfs(i + 1, j)
      if j < len(s2) and s2[j] == s3[i + j]:
          res = res or dfs(i, j + 1)
      return res

  if len(s1) + len(s2) != len(s3):
      return False
  return dfs(0, 0)
```
