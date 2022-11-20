# 前缀和优化 dp dp[i]=dp[left]...+dp[i-1]

- 一般会用到 dp 和 dpSum 数组
- 求出 dp[i]后再更新 dpSum[i]
  [分割数组使得逆序数不大于 k](089%20-%20Partitions%20and%20Inversions%EF%BC%88%E2%98%857%EF%BC%89.py)
  [1997. 访问完所有房间的第一天-前缀和优化](1997.%20%E8%AE%BF%E9%97%AE%E5%AE%8C%E6%89%80%E6%9C%89%E6%88%BF%E9%97%B4%E7%9A%84%E7%AC%AC%E4%B8%80%E5%A4%A9-%E5%89%8D%E7%BC%80%E5%92%8C%E4%BC%98%E5%8C%96.py)

```Python
# dp[i]=dp[lefti]...+dp[i-1] 转化成
dp[i] = dpSum[i - 1] - dpSum[lefti - 1]
```

## 前缀和 dp 的题目特点是某个点的增量/总贡献来源于之前的某段区间

## 滚动数组的前缀和优化写法

1. 只维护 ndp 数组，不维护 ndpSum 数组
   `此时 dpSum 在每个循环开头处理`
   **629. K 个逆序对数组-前缀和优化.-二维**
2. 维护 ndp 数组和 ndpSum 数组
   `注意初始化 先更新dp再更新dpSum`
   `这种写法比较麻烦 容易错`
   **E - Distance Sequence-二维**
   **1977. 划分数字的方案数-index+remain 的前缀和优化 dp**

- 关注边界位置 例如 `取得dp[n]时 preSum[n] 是多少个元素`
  技巧是下边界用 `max(0,cur-nums[i])` 来防止越界
  上边界用 `min(n,cur+nums[i])` 来防止越界
- 关注前缀长度 例如 `preSum[i] - preSum[j] 的 前缀元素个数为 (i - j)`

## 二维 dp 前缀和优化时 需要把连续的那个维度写成数组的第二个维度

例如：前 index 个元素分成 k 组 dp 数组写 dp[count][index] 会更加方便
**6244. 完美分割的方案数-index+remain 的前缀和优化**
