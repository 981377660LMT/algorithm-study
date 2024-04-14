from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

from itertools import accumulate
from math import gcd
from typing import Callable, DefaultDict, List, Optional, Tuple
from collections import defaultdict


def logTrick(
    nums: List[int],
    op: Callable[[int, int], int],
    f: Optional[Callable[[List[Tuple[int, int, int]], int], None]] = None,
) -> DefaultDict[int, int]:
    """
    将 `nums` 的所有非空子数组的元素进行 `op` 操作，返回所有不同的结果和其出现次数.

    Args:
        nums: 1 <= len(nums) <= 1e5.
        op: 与/或/gcd/lcm 中的一种操作，具有单调性.
        f: (interval: List[Tuple[int, int, int]], right: int) -> None
        数组的右端点为right.
        interval 的 leftStart/leftEnd 表示子数组的左端点left的范围.
        interval 的 value 表示该子数组 arr[left,right] 的 op 结果.

    Returns:
        所有不同的结果和其出现次数
    """
    res = defaultdict(int)
    dp = []
    for pos, cur in enumerate(nums):
        for v in dp:
            v[2] = op(v[2], cur)
        dp.append([pos, pos + 1, cur])

        ptr = 0
        for v in dp[1:]:
            if dp[ptr][2] != v[2]:
                ptr += 1
                dp[ptr] = v
            else:
                dp[ptr][1] = v[1]
        dp = dp[: ptr + 1]

        for v in dp:
            res[v[2]] += v[1] - v[0]
        if f is not None:
            f(dp, pos)

    return res


# 给你两个数组 nums 和 andValues，长度分别为 n 和 m。

# 数组的 值 等于该数组的 最后一个 元素。

# 你需要将 nums 划分为 m 个 不相交的连续 子数组，对于第 ith 个子数组 [li, ri]，子数组元素的按位AND运算结果等于 andValues[i]，换句话说，对所有的 1 <= i <= m，nums[li] & nums[li + 1] & ... & nums[ri] == andValues[i] ，其中 & 表示按位AND运算符。


# 返回将 nums 划分为 m 个子数组所能得到的可能的 最小 子数组 值 之和。如果无法完成这样的划分，则返回 -1 。


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def minimumValueSum(self, nums: List[int], andValues: List[int]) -> int:
        def f(interval: List[Tuple[int, int, int]], right: int) -> None:
            for leftStart, leftEnd, value in interval:
                groupByRight[right].append((leftStart, leftEnd, value))

        groupByRight = [[] for _ in range(len(nums))]
        logTrick(nums, lambda x, y: x & y, f)
        # print(groupByRight)
        # for right, v in enumerate(groupByRight):
        #     print(right, v, 666)

        n, m = len(nums), len(andValues)

        @lru_cache(None)
        def dfs(index: int, groupIndex: int) -> int:
            if groupIndex == -1:
                return 0 if index == -1 else INF
            if index == -1:
                return 0 if groupIndex == -1 else INF

            res = INF
            need = andValues[groupIndex]
            groupInfo = groupByRight[index]

            for leftStart, leftEnd, andValue in groupInfo:
                if andValue == need:
                    for i in range(leftStart, leftEnd):
                        cand = dfs(i - 1, groupIndex - 1) + nums[index]
                        res = min2(res, cand)
            return res

        res = dfs(n - 1, m - 1)
        dfs.cache_clear()
        return res if res != INF else -1


class Solution2:
    def minimumValueSum(self, a: List[int], vals: List[int]) -> int:
        n, m = len(a), len(vals)

        @lru_cache(None)
        def dfs(i, k, s):
            if i >= n or k >= m:
                return 0 if i == n and k == m else INF
            ret = dfs(i + 1, k, s & a[i])
            if a[i] & s == vals[k]:
                ret = min(ret, a[i] + dfs(i + 1, k + 1, (1 << 31) - 1))
            return ret

        res = dfs(0, 0, (1 << 31) - 1)
        dfs.cache_clear()
        return res if res != INF else -1


# nums = [1,4,3,3,2], andValues = [0,3,3,2]

if __name__ == "__main__":
    print(Solution().minimumValueSum([1, 4, 3, 3, 2], [0, 3, 3, 2]))
