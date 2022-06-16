**注意`二进制枚举子集`和`枚举子集的子集`的区别**

- 二进制枚举子集：枚举所有二进制数，`O(n*2^n)`
- 状压 dp 枚举所有子集的子集：枚举每个二进制数的子集，`O(3^n)`
  [时间复杂度](https://www.acwing.com/community/content/513/)
  这类题一般都是 `(index1,state2)` 作为 dp 参数

```Python

# 初始化
dp = [False] * (1 << m)
for state in range(1 << m):
    dp[state] = nums[0] >= subsum[state]

# 行间枚举子集的子集转移
for i in range(1, n):
    ndp = [False] * (1 << m)
    for state in range(1 << m):
      # 注意这里可能需要判断g1是空集的情况
      if dp[state]:
          ndp[state] = True
          continue
        g1, g2 = state, 0
        while g1:  # g1是非空子集
            ndp[state] = dp[g2] and nums[i] >= subsum[g1]
            g1 = (g1 - 1) & state
            g2 = state ^ g1

    dp = ndp
```

**枚举所有子集的子集实现方法**

- powerset 函数
- 位运算
