# 前后缀分解模板
from typing import List


class Solution:
    def preSufDp(self, nums: List[int]) -> int:
        def makeDp(nums: List[int]) -> List[int]:
            n = len(nums)
            dp = [0] * (n + 1)
            for i in range(1, n + 1):
                cur = nums[i - 1]
                # your code here
            return dp

        n = len(nums)
        res, preDp, sufDp = 0, makeDp(nums), makeDp(nums[::-1])[::-1]
        for i in range(n):  # 枚举分割点
            res += preDp[i] * sufDp[i + 1]
        return res
