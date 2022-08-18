费用流问题就是要求在所有最大流之中，找到费用最大/最小的问题。

**题目一般是求最优解 最优解肯定是取到最大流时取到**

`二分图最大权匹配可用 KM 算法和 MCMF`
`最小费用流的连续最短路算法复杂度为流量*最短路算法复杂度`

## mcmf 初始化

```Python
V = 150
START, END, OFFSET = 2 * V, 2 * V + 1, V
mcmf = MinCostMaxFlow(2 * v + 2, START, END)
```
