from typing import List
from functools import lru_cache

# 返回每个工人与分配到的自行车之间的曼哈顿距离的最小可能总和。

# 1 <= workers.length <= bikes.length <= 10
# 状压dp


class Solution:
    def assignBikes(self, workers: List[List[int]], bikes: List[List[int]]) -> int:
        dists = [[abs(x - r) + abs(y - c) for r, c in bikes] for x, y in workers]

        @lru_cache(None)
        def dfs(cur: int, visited_bike: int) -> int:
            if cur == len(workers):
                return 0

            res = 0x7FFFFFFF
            for next_bike in range(len(bikes)):
                if ((visited_bike >> next_bike) & 1) == 1:
                    continue
                next_cost = dfs(cur + 1, visited_bike | (1 << next_bike))
                res = min(res, dists[cur][next_bike] + next_cost)

            return res

        res = dfs(0, 0)
        dfs.cache_clear()
        return res


print(Solution().assignBikes([[0, 0], [2, 1]], [[1, 2], [3, 3]]))
print(Solution().assignBikes([[0, 0], [1, 1], [2, 0]], [[1, 0], [2, 2], [2, 1]]))

