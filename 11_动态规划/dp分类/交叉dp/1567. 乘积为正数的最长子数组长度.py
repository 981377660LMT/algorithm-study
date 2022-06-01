# 返回乘积为正数的最长子数组长度。
from typing import List


class Solution:
    def getMaxLen(self, nums: List[int]) -> int:
        res = 0
        pos, neg = 0, 0
        for num in nums:
            # 主意滚动更新相互依赖的要写在一行 不要先更新某个值
            if num > 0:
                pos, neg = pos + 1, neg + 1 if neg != 0 else 0
            elif num < 0:
                pos, neg = neg + 1 if neg != 0 else 0, pos + 1
            else:
                pos, neg = 0, 0
            res = max(res, pos)

        return res

