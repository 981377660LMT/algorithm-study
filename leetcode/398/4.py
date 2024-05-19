from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你有一个 非负 整数 k 。有一个无限长度的台阶，最低 一层编号为 0 。

# 虎老师有一个整数 jump ，一开始值为 0 。虎老师从台阶 1 开始，虎老师可以使用 任意 次操作，目标是到达第 k 级台阶。假设虎老师位于台阶 i ，一次 操作 中，虎老师可以：

# 向下走一级到 i - 1 ，但该操作 不能 连续使用，如果在台阶第 0 级也不能使用。
# 向上走到台阶 i + 2jump 处，然后 jump 变为 jump + 1 。
# 请你返回虎老师到达台阶 k 处的总方案数。


# 注意 ，虎老师可能到达台阶 k 处后，通过一些操作重新回到台阶 k 处，这视为不同的方案。
class Solution:
    def waysToReachStair(self, k: int) -> int:
        @lru_cache(None)
        def dfs(cur: int, preOp: int, jump: int) -> int:
            if cur > k + 1:
                return 0
            res = int(cur == k)
            # down
            if preOp != 0 and cur > 0:
                res += dfs(cur - 1, 0, jump)
            # up
            res += dfs(cur + 2**jump, 1, jump + 1)
            return res

        res = dfs(1, -1, 0)
        dfs.cache_clear()
        return res


print(Solution().waysToReachStair(k=1))
