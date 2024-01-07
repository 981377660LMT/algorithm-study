# 100150. 移除后集合的最多元素数
# https://leetcode.cn/problems/maximum-size-of-a-set-after-removals/description/
# 给你两个下标从 0 开始的整数数组 nums1 和 nums2 ，它们的长度都是偶数 n 。
# 你必须从 nums1 中移除 n / 2 个元素，同时从 nums2 中也移除 n / 2 个元素。
# 移除之后，你将 nums1 和 nums2 中剩下的元素插入到集合 s 中。
# 返回集合 s可能的 最多 包含多少元素。
#
# !删除的优先级为 (当前出现次数、总出现次数)


from typing import List
from heapq import heapify, heappop, heappush
from collections import Counter


class Solution:
    def maximumSetSize(self, nums1: List[int], nums2: List[int]) -> int:
        n = len(nums1)
        counter = Counter(nums1 + nums2)
        counter1 = Counter(nums1)
        counter2 = Counter(nums2)

        pq1 = [(-v, -counter[k], k) for k, v in counter1.items()]  # (freq, removePriority, value)
        heapify(pq1)
        for _ in range(n // 2):
            freq, _, value = heappop(pq1)
            freq = -freq
            counter[value] -= 1
            if counter[value] == 0:
                del counter[value]
            else:
                heappush(pq1, (-(freq - 1), -counter[value], value))

        pq2 = [(-v, -counter[k], k) for k, v in counter2.items()]  # (freq, removePriority, value)
        heapify(pq2)
        for _ in range(n // 2):
            freq, _, value = heappop(pq2)
            freq = -freq
            counter[value] -= 1
            if counter[value] == 0:
                del counter[value]
            else:
                heappush(pq2, (-(freq - 1), -counter[value], value))

        return len(counter)


print(Solution().maximumSetSize([3, 5], [5, 3]))
