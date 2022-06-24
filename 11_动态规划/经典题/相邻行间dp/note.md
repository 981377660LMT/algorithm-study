# 二维状压 dp 的三种境界

1. 普通的相邻行间状态转移 `(row*4^col)`
2. 枚举子集的子集的相邻行间状态转移 `(row*3^col)` (利用了相邻行间同一个列上状态互斥的性质)
3. 轮廓线 dp 的状态转移`(row*col*2^col)`(当前的解只与前 col(+k?) 个格子的状态有关)

## 普通的相邻行间状态转移:

`dfs 或者枚举` 预处理每行可能的状态后，相邻行间进行 dp

优化 1：`可以用一个三进制数表示染色，每位是r/g/b`

```Python
# 在 [0, 3^m) 范围内枚举满足要求的 mask (也可用dfs)
for mask in range(3**m):
    color = list()
    mm = mask
    for i in range(m):
        color.append(mm % 3)
        mm //= 3
    if any(color[i] == color[i + 1] for i in range(m - 1)):
        continue
    valid[mask] = color

```

优化 2：`可以先处理出可能的转移状态邻接表，再进行 dp`

```Python
def colorTheGrid2(self, m: int, n: int) -> int:
        """优化：`可以先处理出可能的转移状态邻接表，再进行 dp`"""
        def dfs(index: int, path: List[int]) -> None:
            if index == n:
                availableStates.append(tuple(path))
                return

            for next in range(3):
                if path and path[-1] == next:
                    continue
                path.append(next)
                dfs(index + 1, path)
                path.pop()

        availableStates: List[State] = []
        dfs(0, [])

        # 邻接表合法转移状态预处理优化 8684 ms => 1512 ms
        adjMap = defaultdict(set)
        for cur in availableStates:
            for next in availableStates:
                if not any(cur[j] == next[j] for j in range(n)):
                    adjMap[cur].add(next)
                    adjMap[next].add(cur)

        dp = [defaultdict(int) for _ in range(m)]
        for state in availableStates:
            dp[0][state] = 1
        for i in range(1, m):
            for preState in dp[i - 1].keys():
                for curState in adjMap[preState]:
                    dp[i][curState] += dp[i - 1][preState]
                    dp[i][curState] %= MOD

        res = 0
        for state in dp[-1].keys():
            res += dp[-1][state]
            res %= MOD
        return res
```

## 枚举子集的子集的相邻行间状态转移

优化 3：`利用相邻行某状态互斥，枚举子集而不是暴力全部枚举，每一行的复杂度从 4^n 降到了 3^n`

三种行间状态转移方式：
https://leetcode-cn.com/problems/EJvmW4/solution/cong-on4m-dao-on3m-by-981377660lmt-kqnb/

1. 预处理每行可能状态+建图+暴力枚举 `2^(2*col)`
   预处理每个行的状态(二进制或者三进制)，每个位置是不贴(0)/单独开(1)/主动联合开(2)，每行产生的花费由状态唯一确定(注意要 cache)，然后相邻行间 dp，看转移是否合法(得到一个有向(带权?)图)，`时间复杂度 O(n*3^2m))，不预处理就是O(n*3^2m*m))`
2. 划分状态，利用相邻行某状态互斥的性质，枚举子集的子集 `3^col`
   `利用相邻行屏障不能重叠的性质把时间复杂度从 O(n*3^2m) 降到了 O(n*3^m)`
   `O(n*3^2m)的行间暴力枚举解法会多计算了两行相同位置都开联合屏障的情况`
   预处理每行的状态+行间枚举子集的子集
   dp[row][state] 定义为第 row 行 某个状态为 state 下的最优解 ，**这个 state 会在相邻行间同一列处互斥**
3. 轮廓线 dp
