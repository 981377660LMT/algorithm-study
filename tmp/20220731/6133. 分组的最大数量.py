"""
你打算将 所有 学生分为一些 有序 的非空分组，其中分组间的顺序满足以下全部条件：
1.第 i 个分组中的学生`总成绩` 小于 第 (i + 1) 个分组中的学生总成绩，对所有组均成立（除了最后一组）。
2.第 i 个分组中的学生`总数 `小于 第 (i + 1) 个分组中的学生总数，对所有组均成立
"""

# !只与长度有关 2满足则1必定满足

from bisect import bisect_right
from typing import List


# 1 3 6
res = [i * (i + 1) // 2 for i in range(1, 10000)]


class Solution:
    def maximumGroups(self, grades: List[int]) -> int:
        return bisect_right(res, len(grades))

        left, right = 0, int(1e9)
        while left <= right:
            mid = (left + right) // 2
            if mid * (mid + 1) // 2 > len(grades):
                right = mid - 1
            else:
                left = mid + 1
        return right


print(Solution().maximumGroups(grades=[10, 6, 12, 7, 3, 5]))
print(Solution().maximumGroups(grades=[8, 8]))
