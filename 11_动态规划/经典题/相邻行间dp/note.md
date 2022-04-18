`dfs 或者枚举` 预处理每行可能的状态后，相邻行间进行 dp

优化 1：`可以先处理出可能的转移状态邻接表，再进行 dp`
优化 2：`可以用一个三进制数表示染色，每位是r/g/b`
优化 3：`枚举子集而不是暴力全部枚举，每一行的复杂度从 2^2n 降到了 3^n`

```Python
# 在 [0, 3^m) 范围内枚举满足要求的 mask
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

        # 优化 8684 ms => 1512 ms
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
