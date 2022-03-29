from typing import List
from functools import lru_cache

# 给你一个整数数组 arr，每一次操作你都可以选择并删除它的一个 回文 子数组 arr[i], arr[i+1], ..., arr[j]（ i <= j）。
# 请你计算并返回从数组中删除所有数字所需的最少操作次数。
# 1 <= arr.length <= 100

# 枚举可能的回文分割点


class Solution:
    def minimumMoves(self, arr: List[int]) -> int:
        @lru_cache(typed=False, maxsize=None)
        def dfs(left: int, right: int) -> int:
            if left >= right:
                return 0
            if right - left <= 2 and arr[left] == arr[right - 1]:
                return 1

            res = 0x3F3F3F3F
            for i in range(left, right):
                # 可能的一段回文
                if arr[i] == arr[left]:
                    # 注意这个max，左边必须删除一个
                    res = min(res, max(1, dfs(left + 1, i)) + dfs(i + 1, right))

            return res

        return dfs(0, len(arr))


print(Solution().minimumMoves(arr=[1, 3, 4, 1, 5]))
# 输出：3
# 解释：先删除 [4]，然后删除 [1,3,1]，最后再删除 [5]。
