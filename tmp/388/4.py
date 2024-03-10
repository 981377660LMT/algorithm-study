from functools import lru_cache
from itertools import accumulate
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个长度为 n 下标从 0 开始的整数数组 nums 和一个 正奇数 整数 k 。

# x 个子数组的能量值定义为 strength = sum[1] * x - sum[2] * (x - 1) + sum[3] * (x - 2) - sum[4] * (x - 3) + ... + sum[x] * 1 ，其中 sum[i] 是第 i 个子数组的和。更正式的，能量值是满足 1 <= i <= x 的所有 i 对应的 (-1)i+1 * sum[i] * (x - i + 1) 之和。

# 你需要在 nums 中选择 k 个 不相交子数组 ，使得 能量值最大 。

# 请你返回可以得到的 最大能量值 。


# 注意，选出来的所有子数组 不 需要覆盖整个数组。
def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def maximumStrength(self, nums: List[int], k: int) -> int:
        @lru_cache(None)
        def dfs(index: int, remain: int) -> int:
            if remain < 0:
                return -INF
            if index == n:
                return 0 if remain == 0 else -INF

            minus = 1 if ((k - remain) % 2 == 0) else -1
            res = dfs(index + 1, remain) + minus * nums[index] * remain
            if remain > 0:
                tmp = dfs(index + 1, remain - 1) + minus * nums[index] * remain
                res = max2(res, tmp)
            return res

        n = len(nums)
        res = -INF
        for i in range(n):
            res = max2(res, dfs(i, k))
        dfs.cache_clear()
        return res


# nums = [1,2,3,-1,2], k = 3
# nums = [-1,-2,-3], k = 1


print(Solution().maximumStrength([-1, -2, -3], 2))
# [-99,85]
# 1
print(Solution().maximumStrength([-99, 85], 1))
