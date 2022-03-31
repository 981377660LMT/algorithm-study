# 删除连续字符,使源字符串没有连续字符
# 求最小花费
from itertools import groupby
from typing import List


class Solution:
    def solve(self, s: str, nums: List[int]):
        cost = sum(nums)
        for _, group in groupby(zip(s, nums), key=lambda pair: pair[0]):
            cost -= max(c for _, c in group)
        return cost


print(Solution().solve("aacd", [1, 2, 3, 4]))
