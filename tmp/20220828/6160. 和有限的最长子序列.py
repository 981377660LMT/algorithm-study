from bisect import bisect_right
from itertools import accumulate
from typing import List


class Solution:
    def answerQueries(self, nums: List[int], queries: List[int]) -> List[int]:
        """
        返回一个长度为 m 的数组 answer ,
        其中 answer[i] 是 nums 中 元素之和小于等于 queries[i] 的 子序列 的 最大 长度  。

        子序列排序
        """
        preSum = [0] + list(accumulate(sorted(nums)))
        return [bisect_right(preSum, q) - 1 for q in queries]

    # !把子序列改成子数组要怎么做？ 滑动窗口
