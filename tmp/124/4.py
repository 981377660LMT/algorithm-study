from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始只包含 正 整数的数组 nums 。

# 一开始，你可以将数组中 任意数量 元素增加 至多 1 。

# 修改后，你可以从最终数组中选择 一个或者更多 元素，并确保这些元素升序排序后是 连续 的。比方说，[3, 4, 5] 是连续的，但是 [3, 4, 6] 和 [1, 1, 2, 3] 不是连续的。
# 请你返回 最多 可以选出的元素数目。


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def maxSelectedElements(self, nums: List[int]) -> int:
        @lru_cache(None)
        def dfs(cur: int, borrow: bool) -> int:
            remain = counter[cur - 1] - borrow
            if remain < 0:
                return 0
            if remain > 0:
                return 1 + dfs(cur - 1, False)
            if counter[cur - 2] <= 0:
                return 1
            return 1 + dfs(cur - 1, True)

        counter = defaultdict(int)
        max_ = 0
        for v in nums:
            counter[v] += 1
            max_ = max(max_, v)
        res = 1
        for m in range(1, max_ + 2):
            if counter[m] == 0 and counter[m - 1] == 0:
                continue
            res = max2(res, dfs(m, counter[m] == 0))
        dfs.cache_clear()
        return res


# nums = [2,1,5,1,1]

print(Solution().maxSelectedElements([2, 1, 5, 1, 1]))
# [8,13,18,10,16,19,11,17,15,18,9,12,15,8,9,14,7]
print(
    Solution().maxSelectedElements([8, 13, 18, 10, 16, 19, 11, 17, 15, 18, 9, 12, 15, 8, 9, 14, 7])
)
