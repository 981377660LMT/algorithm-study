from heapq import heappop, heappush
from typing import List

# 给你两个下标从 0 开始的整数数组 nums1 和 nums2 ，两者长度都是 n ，
# 再给你一个正整数 k 。你必须从 nums1 中选一个长度为 k 的 子序列 对应的下标。

# 对于选择的下标 i0 ，i1 ，...， ik - 1 ，你的 分数 定义如下：
# !nums1 中下标对应元素求和，乘以 nums2 中下标对应元素的 最小值 。
# 用公示表示： (nums1[i0] + nums1[i1] +...+ nums1[ik - 1]) * min(nums2[i0] , nums2[i1], ... ,nums2[ik - 1]) 。
# 请你返回 最大 可能的分数。

# 一个数组的 子序列 下标是集合 {0, 1, ..., n-1} 中删除若干元素得到的剩余集合，也可以不删除任何元素。

# 解:
# !枚举最小值，然后用堆维护最大的k个数


class Solution:
    def maxScore(self, nums1: List[int], nums2: List[int], k: int) -> int:
        pq = []
        res = 0
        curSum = 0
        for min_, add in sorted(zip(nums2, nums1), reverse=True):
            heappush(pq, add)
            curSum += add
            if len(pq) > k:
                curSum -= heappop(pq)
            if len(pq) == k:
                res = max(res, curSum * min_)

        return res


assert Solution().maxScore(nums1=[1, 3, 3, 2], nums2=[2, 1, 3, 4], k=3) == 12
