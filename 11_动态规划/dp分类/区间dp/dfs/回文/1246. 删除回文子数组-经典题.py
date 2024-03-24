from typing import List
from functools import lru_cache

# 给你一个整数数组 arr，每一次操作你都可以选择并删除它的一个 回文 子数组 arr[i], arr[i+1], ..., arr[j]（ i <= j）。
# 请你计算并返回从数组中删除所有数字所需的最少操作次数。

# !注意删除完之后右侧元素会向左移动
# 1 <= arr.length <= 100

# 枚举可能的回文分割点

# 1 <= arr.length <= 100
# 1 <= arr[i] <= 20


class Solution:
    def minimumMoves(self, arr: List[int]) -> int:
        @lru_cache(None)
        def dfs(left: int, right: int) -> int:
            """删除[left,right]回文子数组的最少操作次数"""
            if left > right:
                return 0
            if right - left <= 1:
                return 1 if arr[left] == arr[right] else 2

            res = int(1e18)
            if arr[left] == arr[right]:
                res = min(res, dfs(left + 1, right - 1))  # 追加到两侧 不增加操作次数

            for mid in range(left, right):  # 左端点必须找到另个一相等的数配对消除 枚举分割点 分成左右区间
                # 可能的一段回文
                if arr[mid] == arr[left]:
                    res = min(res, dfs(left, mid) + dfs(mid + 1, right))
            return res

        n = len(arr)
        return dfs(0, n - 1)


print(Solution().minimumMoves(arr=[1, 2]))
print(Solution().minimumMoves(arr=[1, 3, 4, 1, 5]))
print(Solution().minimumMoves(arr=[1, 14, 18, 20, 14]))  # 3
# 输出：3
# 解释：先删除 [4]，然后删除 [1,3,1]，最后再删除 [5]。
print(Solution().minimumMoves(arr=[1, 4, 1, 1, 2, 3, 2, 1]))
