# 返回 两个数组中 公共的 、长度最长的子数组的长度 。

from typing import List

# 1923. 最长公共子路径 多个数组的最长公共子数组


class Solution:
    def findLength(self, nums1: List[int], nums2: List[int]) -> int:
        """注意不是最长公共子序列 而是最长公共子数组
        
        滑窗枚举对齐起点
        时间复杂度O((n+m)*min(n,m))
        """

        def cal(offset1: int, offset2: int, length: int) -> int:
            res = dp = 0
            for i in range(length):
                if nums1[offset1 + i] == nums2[offset2 + i]:
                    dp += 1
                    res = max(res, dp)
                else:
                    dp = 0
            return res

        n1, n2 = len(nums1), len(nums2)
        res = 0
        for start in range(n1):
            length = min(n1 - start, n2)
            res = max(res, cal(start, 0, length))

        for start in range(n2):
            length = min(n2 - start, n1)
            res = max(res, cal(0, start, length))
        return res

