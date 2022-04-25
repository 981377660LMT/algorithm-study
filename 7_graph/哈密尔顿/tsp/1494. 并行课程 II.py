from typing import List
from functools import lru_cache
from itertools import combinations

# 在一个学期中，你 最多 可以同时上 k 门课，前提是这些课的先修课在之前的学期里已经上过了。
# 请你返回上完所有课最少需要多少个学期。题目保证一定存在一种上完所有课的方式。

# 1 <= n <= 15  暗示2^n 状压

# 这道题其实是说，n个用时相同的任务（work）有先后依赖关系，现只有k台机器（parallelism 并行度），最少用时多少能完成。最大深度为depth。


class Solution:
    def minNumberOfSemesters(self, n: int, dependencies: List[List[int]], k: int) -> int:
        @lru_cache(None)
        def dfs(state: int) -> int:
            if state == target:
                return 0

            res = int(1e20)
            nexts = [i for i in range(n) if not state & (1 << i) and deps[i] & state == deps[i]]
            for sub in combinations(nexts, min(k, len(nexts))):
                res = min(res, 1 + dfs(state | sum(1 << i for i in sub)))
            return res

        deps = [0] * n
        for pre, cur in dependencies:
            pre, cur = pre - 1, cur - 1
            deps[cur] |= 1 << pre

        target = (1 << n) - 1
        res = dfs(0)
        dfs.cache_clear()
        return res


print(Solution().minNumberOfSemesters(n=4, dependencies=[[2, 1], [3, 1], [1, 4]], k=2))
# 在第一个学期中，我们可以上课程 2 和课程 3 。然后第二个学期上课程 1 ，第三个学期上课程 4 。
print(Solution().minNumberOfSemesters(n=5, dependencies=[[2, 1], [3, 1], [4, 1], [1, 5]], k=2))
