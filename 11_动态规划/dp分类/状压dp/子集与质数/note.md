这类状压的特点是数字大小特别小(nums[i]<=30)
所以要从值域的角度去考虑

```python
dp[index][state] 表示当前在值域上遍历到了keys[index],素数集合状态为state的方案数

@lru_cache(None)
def dfs(index:int,state:int)->int:
  ...
```

- 30 以内的质数有 10 个 [2,3,5,7,11,13,17,19,23,29]
