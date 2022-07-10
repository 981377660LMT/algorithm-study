from typing import List
from functools import lru_cache

# 返回每个工人与分配到的自行车之间的曼哈顿距离的最小可能总和。

# 1 <= workers.length <= bikes.length <= 10
# 状压dp


class Solution:
    def assignBikes(self, workers: List[List[int]], bikes: List[List[int]]) -> int:
        @lru_cache(None)
        def dfs(index: int, visited: int) -> int:
            if index == len(workers):
                return 0

            res = int(1e20)
            for j in range(len(bikes)):
                if (visited >> j) & 1:
                    continue
                res = min(res, dists[index][j] + dfs(index + 1, visited | (1 << j)))

            return res

        dists = [[abs(x - r) + abs(y - c) for r, c in bikes] for x, y in workers]
        res = dfs(0, 0)
        dfs.cache_clear()
        return res


print(Solution().assignBikes([[0, 0], [2, 1]], [[1, 2], [3, 3]]))
print(Solution().assignBikes([[0, 0], [1, 1], [2, 0]], [[1, 0], [2, 2], [2, 1]]))
