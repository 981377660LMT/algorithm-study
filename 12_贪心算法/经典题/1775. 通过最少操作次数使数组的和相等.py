from typing import List
from collections import Counter

# 1775. 通过最少操作次数使数组的和相等
# 两个数组中的所有值都在 1 到 6 之间（包含 1 和 6）。
# 每次操作中，你可以选择 任意 数组中的任意一个整数，
# 将它变成 1 到 6 之间 任意 的值（包含 1 和 6）。
# 请你返回使 nums1 中所有数的和与 nums2 中所有数的和相等的最少操作次数。
# 如果无法使两个数组的和相等，请返回 -1 。


class Solution:
    def minOperations(self, nums1: List[int], nums2: List[int]) -> int:
        if 6 * len(nums1) < len(nums2) or 6 * len(nums2) < len(nums1):
            return -1  # impossible

        sum1, sum2 = sum(nums1), sum(nums2)
        if sum1 < sum2:
            nums1, nums2 = nums2, nums1
            sum1, sum2 = sum2, sum1

        diff = sum1 - sum2
        freq = Counter(num - 1 for num in nums1) + Counter(6 - num for num in nums2)

        res = 0
        for delta in range(5, 0, -1):
            for _ in range(freq[delta]):
                if diff <= 0:
                    break
                res += 1
                diff -= delta
        return res


print(Solution().minOperations(nums1=[1, 2, 3, 4, 5, 6], nums2=[1, 1, 2, 2, 2, 2]))
# 输出：3
# 解释：你可以通过 3 次操作使 nums1 中所有数的和与 nums2 中所有数的和相等。以下数组下标都从 0 开始。
# - 将 nums2[0] 变为 6 。 nums1 = [1,2,3,4,5,6], nums2 = [6,1,2,2,2,2] 。
# - 将 nums1[5] 变为 1 。 nums1 = [1,2,3,4,5,1], nums2 = [6,1,2,2,2,2] 。
# - 将 nums1[2] 变为 2 。 nums1 = [1,2,2,4,5,1], nums2 = [6,1,2,2,2,2] 。
