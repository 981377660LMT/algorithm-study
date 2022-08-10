**注意`二进制枚举子集`和`枚举子集的子集`的区别**

- 二进制枚举子集：枚举所有二进制数，`O(n*2^n)`
- 状压 dp 枚举所有子集的子集：枚举每个二进制数的子集，`O(3^n)`
  [时间复杂度](https://www.acwing.com/community/content/513/)
  这类题一般都是 `(index1,state2)` 作为 dp 参数

```Python

# 初始化所有状态
dp = [False] * (1 << m)
for state in range(1 << m):
    dp[state] = nums[0] >= subsum[state]

# 行间枚举子集的子集转移 O(n*3^m)
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

## 枚举所有子集的子集实现方法

- powerset 函数
- 位运算

## 需要枚举子集的子集的特征:

1. 解决元素的`分组`问题 即哪几个元素分到一组会取到最优解

   - 小朋友分饼干、顾客分数字、工人分配任务、点分组 (两个维度的分组)
   - 直线分组、兔子分组 (一个维度内部分组)
     `dp参数定义为 (index,state) 复杂度O(k*3^n)` or
     `dp参数定义为 (state) 复杂度O(3^n)`

2. n 很小，一般 n<=16
