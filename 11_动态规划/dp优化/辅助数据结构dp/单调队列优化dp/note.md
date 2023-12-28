单调队列优化 DP
`dequeue 维护区间最值`

1. 在连续 mid+1 个烽火台中至少要有一个发出信号
   [](1089.%20%E7%83%BD%E7%81%AB%E4%BC%A0%E9%80%92.py)
   [](1090.%20%E7%BB%BF%E8%89%B2%E9%80%9A%E9%81%93-%E5%8D%95%E8%B0%83%E9%98%9F%E5%88%97+%E4%BA%8C%E5%88%86.py)
2. 环形单调队列
   [](1088.%20%E6%97%85%E8%A1%8C%E9%97%AE%E9%A2%98-%E6%BB%91%E5%8A%A8%E7%AA%97%E5%8F%A3%E6%9C%80%E5%B0%8F%E5%80%BC.py)

---

如果 dp 转移中出现`滑动窗口最值`，就可以用单调队列优化

dp[i] = min(dp[j] for j in [i-k, i-1]) + C[i]

---

`用 n+1 还是 n  长度的 dp 数组取决于实际情况`

```python
# dp[i] 表示前 i 个元素中，以第 i 个元素结尾的子序列元素和的最大值(0<=i<n)
# dp[i] = max(dp[i], max(dp[i - k] ,..., dp[i-1], 0) + nums[i])
# res = max(dp)
dp = [-INF] * n
seg = MonoQueue[Tuple[int, int]](lambda x, y: x[0] > y[0])
for i, num in enumerate(nums):
  while seg and i - seg.head()[1] > k:  # 1.不在窗口内的元素出队
      seg.popleft()
  preMax = max(0, seg.head()[0]) if seg else 0
  dp[i] = max(dp[i], preMax + num)  # 2.更新dp
  seg.append((dp[i], i))  # 3.入队
return max(dp)

# dp[i]:获得前i个水果`且最后一个水果i是花钱买的`，所需的最少金币数(0<=i<=n)
# dp[0]=0
# dp[i]=min(dp[j]+prices[i-1]) | i//2<=j<i)
# !维护窗口内的`dp[j]`的最小值即可
# !答案为min(dp[(n+1)//2],dp[n//2+1],...,dp[n])
queue = MonoQueue[Tuple[int, int]](lambda x, y: x[0] < y[0])  # (value, index)
dp = [INF] * (n + 1)
dp[0] = 0
queue.append((dp[0], 0))
for i in range(1, n + 1):
   while queue and queue.head()[1] < i // 2:
       queue.popleft()
   preMin = queue.head()[0] if queue else INF
   dp[i] = min(dp[i], preMin + prices[i - 1])
   queue.append((dp[i], i))
return min(dp[(n + 1) // 2 : n + 1])
```

---

实际上，优化 dp 时线段树写出来可读性更好，建议使用线段树代替
