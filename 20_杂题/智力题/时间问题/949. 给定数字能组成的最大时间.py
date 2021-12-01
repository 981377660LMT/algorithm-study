from typing import List
from itertools import permutations


class Solution:
    def largestTimeFromDigits(self, arr: List[int]) -> str:
        return max(
            ("{}{}:{}{}".format(*p) for p in permutations(arr) if p[:2] < (2, 4) and p[2] < 6),
            default='',
        )


print(Solution().largestTimeFromDigits([1, 2, 3, 4]))
# 输出："23:41"
# 解释：有效的 24 小时制时间是 "12:34"，"12:43"，"13:24"，"13:42"，"14:23"，"14:32"，"21:34"，"21:43"，"23:14" 和 "23:41" 。这些时间中，"23:41" 是最大时间。

