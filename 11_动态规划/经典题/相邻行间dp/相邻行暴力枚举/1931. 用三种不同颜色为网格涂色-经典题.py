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
        for state in availableStates:
            dp[0][state] = 1
        for i in range(1, m):
            for state in availableStates:
                for preState in dp[i - 1].keys():
                    if not any(preState[j] == state[j] for j in range(n)):
                        dp[i][state] += dp[i - 1][preState]
                        dp[i][state] %= MOD

        res = 0
        for state in dp[-1].keys():
            res += dp[-1][state]
            res %= MOD
        return res

    def colorTheGrid2(self, m: int, n: int) -> int:
        """优化：`可以先处理出可能的转移状态邻接表，再进行 dp`"""
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


print(Solution().colorTheGrid1(m=1, n=2))
print(Solution().colorTheGrid2(m=1, n=1))
