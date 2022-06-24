from collections import defaultdict
from typing import List, Tuple

MOD = int(1e9 + 7)
State = Tuple[int, ...]

# 1 <= m <= 5
# 1 <= n <= 1000


class Solution:
    def colorTheGrid1(self, m: int, n: int) -> int:
        """思路同2184，先状压处理每行的可能状态，再dp处理相邻行间的状态"""
        n, m = m, n

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

        dp = [defaultdict(int) for _ in range(m)]
        for cur in availableStates:
            dp[0][cur] = 1
        for i in range(1, m):
            for cur in availableStates:
                for pre in dp[i - 1]:
                    if not any(pre[j] == cur[j] for j in range(n)):
                        dp[i][cur] += dp[i - 1][pre]
                        dp[i][cur] %= MOD

        res = 0
        for cur in dp[-1]:
            res += dp[-1][cur]
            res %= MOD
        return res

    def colorTheGrid2(self, m: int, n: int) -> int:
        """优化：`可以先处理出可能的转移状态邻接表，再进行 dp`"""

        def dfs(index: int, path: List[int]) -> None:
            if index == COL:
                allStates.append(tuple(path))
                return

            for next in range(3):
                if path and path[-1] == next:
                    continue
                path.append(next)
                dfs(index + 1, path)
                path.pop()

        ROW, COL = sorted([m, n], reverse=True)

        allStates: List[State] = []
        dfs(0, [])

        # 优化 8684 ms => 1512 ms
        adjMap = defaultdict(set)
        for cur in allStates:
            for next in allStates:
                if not any(cur[j] == next[j] for j in range(COL)):
                    adjMap[cur].add(next)
                    adjMap[next].add(cur)

        dp = defaultdict(int, {s: 1 for s in allStates})

        for _ in range(1, ROW):
            ndp = defaultdict(int)
            for pre in dp:
                for cur in adjMap[pre]:
                    ndp[cur] += dp[pre]
                    ndp[cur] %= MOD
            dp = ndp

        res = 0
        for count in dp.values():
            res += count
            res %= MOD
        return res


print(Solution().colorTheGrid1(m=1, n=2))
print(Solution().colorTheGrid2(m=1, n=1))
