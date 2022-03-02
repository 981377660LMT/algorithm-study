# 你需要用 红，黄，绿 三种颜色之一给每一个格子上色，且确保`相邻格子颜色不同`（
# 给你网格图的行数 n 。

# 1 <= n <= 5000


# 总结：Two pattern for each row, 121 and 123.
# We consider the state of the first row,
# pattern 121: 121, 131, 212, 232, 313, 323.
# pattern 123: 123, 132, 213, 231, 312, 321.
# So we initialize a121 = 6, a123 = 6.

# We consider the next possible for each pattern:
# Patter 121 can be followed by: 212, 213, 232, 312, 313
# Patter 123 can be followed by: 212, 231, 312, 232

# 121 => three 121, two 123
# 123 => two 121, two 123

# So we can write this dynamic programming transform equation:
# b121 = a121 * 3 + a123 * 2
# b123 = a121 * 2 + a123 * 2


from collections import defaultdict
from typing import List, Tuple

MOD = int(1e9 + 7)
State = Tuple[int, int, int]


class Solution:
    def numOfWays(self, n: int) -> int:
        """思路同2184，先状压处理每行的可能状态，再dp处理相邻行间的状态"""

        def dfs(index: int, path: List[int]) -> None:
            if index == 3:
                availableStates.append(tuple(path))
                return
            for next in (1, 2, 3):
                if path and path[-1] == next:
                    continue
                path.append(next)
                dfs(index + 1, path)
                path.pop()

        availableStates: List[State] = []
        dfs(0, [])

        dp = [defaultdict(int) for _ in range(n)]
        for state in availableStates:
            dp[0][state] = 1
        for i in range(1, n):
            for state in availableStates:
                for preState in dp[i - 1].keys():
                    if not any(preState[j] == state[j] for j in range(3)):
                        dp[i][state] += dp[i - 1][preState]
                        dp[i][state] %= MOD

        res = 0
        for state in dp[-1].keys():
            res += dp[-1][state]
            res %= MOD
        return res


print(Solution().numOfWays(n=2))
# 54
