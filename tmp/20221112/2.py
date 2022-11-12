from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你整数 zero ，one ，low 和 high ，我们从空字符串开始构造一个字符串，每一步执行下面操作中的一种：

# 将 '0' 在字符串末尾添加 zero  次。
# 将 '1' 在字符串末尾添加 one 次。
# 以上操作可以执行任意次。

# 如果通过以上过程得到一个 长度 在 low 和 high 之间（包含上下边界）的字符串，那么这个字符串我们称为 好 字符串。

# 请你返回满足以上要求的 不同 好字符串数目。由于答案可能很大，请将结果对 109 + 7 取余 后返回。


class Solution:
    def countGoodStrings(self, low: int, high: int, zero: int, one: int) -> int:
        def cal(upper: int) -> int:
            @lru_cache(None)
            def dfs(curLength: int) -> int:
                if curLength >= upper:
                    return 1 if curLength == upper else 0
                return (1 + dfs(curLength + zero) + dfs(curLength + one)) % MOD

            return dfs(0)

        return (cal(high) - cal(low - 1)) % MOD


print(Solution().countGoodStrings(low=3, high=3, zero=1, one=1))
