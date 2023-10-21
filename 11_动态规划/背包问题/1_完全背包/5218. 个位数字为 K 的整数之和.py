from functools import lru_cache
from typing import List, Tuple
from collections import defaultdict, Counter


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个整数 num 和 k ，考虑具有以下属性的正整数多重集：

# 每个整数个位数字都是 k 。
# 所有整数之和是 num 。


# 返回该多重集的最小大小，如果不存在这样的多重集，返回 -1 。

# 0 <= num <= 3000
# 0 <= k <= 9
class Solution:
    def minimumNumbers(self, num: int, k: int) -> int:
        """完全背包"""

        @lru_cache(None)
        def dfs(remain: int) -> int:
            if remain < 0:
                return INF
            if remain == 0:
                return 0

            res = INF
            for select in nums:
                if remain - select < 0:
                    break
                res = min(res, 1 + dfs(remain - select))
            return res

        nums = []
        for v in range(k, num + 1, 10):
            if v != 0:
                nums.append(v)

        res = dfs(num)
        dfs.cache_clear()
        return res if res < int(1e10) else -1

    def minimumNumbers2(self, num: int, k: int) -> int:
        """枚举+数学

        若多重集里有 n 个整数，那么这些整数之和为 (10*∑ai + n*k)
        只要 (num - nk) 能被 10 整除且大等于 0,就存在一个大小为 n 的集合。返回最小的符合要求的 n 即可
        """

        if num == 0:
            return 0
        for i in range(1, num + 1):
            diff = num - i * k
            if diff % 10 == 0 and diff >= 0:
                return i
        return -1


print(Solution().minimumNumbers(num=4, k=0))
print(Solution().minimumNumbers(num=58, k=9))
