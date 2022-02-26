https://leetcode-cn.com/problems/shortest-path-visiting-all-nodes/solution/gtalgorithm-tu-jie-fa-ba-hardbian-cheng-v5knb/
来源于数学家哈密尔顿与十八面体(ps:只有五种正多面体，4,6,8,12,20)
是否存在哈密尔顿回路是一个 NP 难问题

状压+记忆化搜索(memo)解决哈密尔顿/tsp 问题
旅行商问题（TSP）：给定一系列城市和每对城市之间的距离，求解访问每一座城市`一次`并回到起始城市的最短回路

```TS
declare function dfs(cur:number,visited:number,cost:number):boolean
```

dfs 需要记忆化

不用记忆化暴力回溯:`O(n!)` n<=10 量级
使用记忆化搜索:`(O(n*2^n))` n<=20 量级 (`等差数列*等比数列`的复杂度)
mask 这一维度因此可以看成是进行了 2^n 次常规的广度优先搜索
最坏情况下为完全图，有边数 = O(n^2)

所以 n=15 一般是暗示记忆化搜索

## 模板:

bfs`864. 获取所有钥匙的最短路径.py`
dfs`847. 访问所有节点的最短路径.ts`
关键点:

1. bfs:三个参数存队列,dfs:三个参数为传参
2. visited 记录(数组或者 set 记录(`cur,visitedState`))
3. `visited === targe`t 则返回 cost
