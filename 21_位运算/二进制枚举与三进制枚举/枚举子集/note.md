https://leetcode-cn.com/problems/subsets/solution/zi-ji-by-leetcode-solution/

**枚举子集，时间复杂度`O(n*2^n)`**

- 枚举 state + check

- dfs (yield 写法)
  `O(n*2^n)，这里的每种状态需要 O(n) 的时间来构造子集包含了push和pop操作`

- powerset 函数(通用)

- 滚动集合更新

**求每个子集的和(某个相关的值)，可用 dp/dfs 优化到 O(2^n)**

- dp
- dfs(比较好，可以剪枝)
- 滚动集合更新
