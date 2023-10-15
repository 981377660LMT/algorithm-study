我们知道滑动窗口适合在题目要求**连续**的情况下使用， 而前缀和也是如此。
前缀和，这个概念其实很容易理解，即一个数组中，第 n 位存储的是数组前 n 个数字的和。
我们可以使用公式 pre[𝑖]=pre[𝑖−1]+nums[𝑖]得到每一位前缀和的值，从而通过前缀和进行相应的计算和解题。
When seeing the word **consecutive sequence**, think about using **Prefix Sum** strategy.

两个 atMostK:
**不超过 k 种元素的子数组个数**:水果成栏问题,K 个不同整数的子数组问题：滑窗
**全部元素都不大于 k 的子数组个数**：`795. 区间子数组个数`

```Python
preSum = [0]
for num in nums:
    preSum.append(preSum[-1] + num)
```

此时：
`注意 preSum[i+1]的含义:arr[0]+arr[1]+...+arr[i]`不要弄混了
**i+1 表示前 i+1 个数 因此是 0-i 这 i+1 项**
`[i,j]这一段的和是 preSum[j+1]-preSum[i]` 一共`(j-i+1)`个数

确定前缀和下标的两点:

1. 和切片一样，两个 index 相减等于长度
2. 带入特殊值验证 (例如长为整个数组长度 n 时对应的索引)

---

静态的区间`计数`查询可以考虑前缀和
