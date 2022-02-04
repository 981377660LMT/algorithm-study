from collections import defaultdict
from functools import lru_cache
from typing import List


class Point:
    def __init__(self, a=0, b=0):
        self.x = a
        self.y = b


# 1<n≤20
# 请问牛牛能组成的新集合多少种。
# 时间复杂度:O(2^n)
class Solution:
    def solve(self, n: int, m: int, limit: List[Point]) -> int:
        def dfs(cur: int, visited: int) -> int:
            if cur == n + 1:
                return 1
            res = dfs(cur + 1, visited)
            if not visited & inValid[cur]:
                res += dfs(cur + 1, visited | (1 << cur))
            return res

        # write code here
        # 对每个数字，`预处理`它不能和哪些数在一起
        inValid = defaultdict(int)
        for point in limit:
            x, y = point.x, point.y
            inValid[x] |= 1 << y
            inValid[y] |= 1 << x

        return dfs(1, 0)


print(Solution().solve(3, 2, [Point(1, 2), Point(2, 3)]))
