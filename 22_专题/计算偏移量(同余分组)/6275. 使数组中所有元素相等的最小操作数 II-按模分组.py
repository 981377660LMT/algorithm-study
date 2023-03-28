# 6275. 使数组中所有元素相等的最小操作数 II-按模分组
from typing import List
from collections import defaultdict, Counter

# 给你两个整数数组 nums1 和 nums2 ，两个数组长度都是 n ，
# 再给你一个整数 k 。你可以对数组 nums1 进行以下操作：

# 选择两个下标 i 和 j ，将 nums1[i] 增加 k ，将 nums1[j] 减少 k 。
# !换言之，nums1[i] = nums1[i] + k 且 nums1[j] = nums1[j] - k 。
# 如果对于所有满足 0 <= i < n 都有 num1[i] == nums2[i] ，那么我们称 nums1 等于 nums2 。

# 请你返回使 nums1 等于 nums2 的 最少 操作数。如果没办法让它们相等，请你返回 -1 。

# !0 <= k <= 1e5 注意k可能为0


class Solution:
    def minOperations(self, nums1: List[int], nums2: List[int], k: int) -> int:
        """对应下标需要相等"""
        if sum(nums1) != sum(nums2):
            return -1
        if k == 0:
            return 0 if nums1 == nums2 else -1

        res = 0
        for a, b in zip(nums1, nums2):
            div, mod = divmod(a - b, k)
            if mod != 0:
                return -1
            res += abs(div)
        return res // 2

    def minOperations2(self, nums1: List[int], nums2: List[int], k: int) -> int:
        """不看下标,counter相等,需要按mod k分组后排序,每个组内对应相等"""
        if sum(nums1) != sum(nums2):
            return -1
        if k == 0:
            return 0 if Counter(nums1) == Counter(nums2) else -1

        n = len(nums1)
        g1, g2 = defaultdict(list), defaultdict(list)
        for i in range(n):
            g1[nums1[i] % k].append(nums1[i])
            g2[nums2[i] % k].append(nums2[i])
        for g in g1.values():
            g.sort()
        for g in g2.values():
            g.sort()

        res = 0
        for mod in range(k):
            if len(g1[mod]) != len(g2[mod]):
                return -1
            for a, b in zip(g1[mod], g2[mod]):
                res += abs(a - b) // k
        return res // 2


print(Solution().minOperations(nums1=[4, 3, 1, 4], nums2=[1, 3, 7, 1], k=3))
print(Solution().minOperations2(nums1=[4, 3, 1, 4], nums2=[1, 3, 7, 1], k=3))
