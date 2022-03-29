from typing import List, Tuple
from functools import lru_cache

# 1 <= n <= 14
# max(tasks[i]) <= sessionTime <= 15

# 第 i 个任务需要花费 tasks[i] 小时完成。一个 工作时间段 中，你可以 至多 连续工作 sessionTime 个小时，然后休息一会儿。
# 返回完成所有任务所需要的 最少 数目的 工作时间段 。
# 加了限制条件后，就变成了贪心 881. 救生艇

# 总结:
# 技巧：dfs返回多个值

# Time complexity is O(2^n * n), because we have 2^n masks and O(n) transitions from given mask. Space complexity is O(2^n).

INF = 0x7FFFFFFF


class Solution:
    def minSessions(self, tasks: List[int], sessionTime: int) -> int:
        @lru_cache(None)
        def dfs(state: int, remain: int) -> int:
            if state == target:
                return 1

            res = INF
            for i in range(n):
                if state & (1 << i):
                    continue

                nextState = state | (1 << i)
                nextCost = tasks[i]
                nextRemain = remain - nextCost if remain >= nextCost else sessionTime - nextCost
                res = min(res, int(nextCost > remain) + dfs(nextState, nextRemain))
            return res

        tasks.sort(reverse=True)
        n = len(tasks)
        target = (1 << n) - 1

        res = dfs(0, sessionTime)
        dfs.cache_clear()
        return res


print(Solution().minSessions(tasks=[3, 1, 3, 1, 1], sessionTime=8))
# 输出：2
# 解释：你可以在两个工作时间段内完成所有任务。
# - 第一个工作时间段：完成除了最后一个任务以外的所有任务，花费 3 + 1 + 3 + 1 = 8 小时。
# - 第二个工作时间段，完成最后一个任务，花费 1 小时。

