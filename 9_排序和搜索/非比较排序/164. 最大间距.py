# 给定一个无序的数组 nums，返回 数组在排序之后，相邻元素之间最大的差值 。如果数组元素个数小于 2，则返回 0 。
# 您必须编写一个在「线性时间」内运行并使用「线性额外空间」的算法。


from typing import List


INF = int(1e20)


class Solution:
    def maximumGap(self, nums: List[int]) -> int:
        if len(nums) < 2:
            return 0
        min_, max_ = min(nums), max(nums)
        n = len(nums)
        size = max(1, (max_ - min_) // (n - 1))
        buckets = [[INF, -INF] for _ in range((max_ - min_) // size + 1)]
        for v in nums:
            b = buckets[(v - min_) // size]
            b[0] = min(b[0], v)
            b[1] = max(b[1], v)

        preMax = INF
        res = 0
        for mn, mx in buckets:
            if mn == INF:
                continue
            # 桶内最小值，减去上一个非空桶的最大值
            res = max(res, mn - preMax)
            preMax = mx

        return res
