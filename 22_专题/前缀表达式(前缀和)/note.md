我们知道滑动窗口适合在题目要求**连续**的情况下使用， 而前缀和也是如此。
前缀和，这个概念其实很容易理解，即一个数组中，第 n 位存储的是数组前 n 个数字的和。
我们可以使用公式 pre[𝑖]=pre[𝑖−1]+nums[𝑖]得到每一位前缀和的值，从而通过前缀和进行相应的计算和解题。
When seeing the word **consecutive sequence**, think about using **Prefix Sum** strategy.

两个 atMostK:
**不超过 k 种元素的子数组个数**:水果成栏问题,K 个不同整数的子数组问题
**全部元素都不大于 k 的子数组个数**：区间子数组个数
