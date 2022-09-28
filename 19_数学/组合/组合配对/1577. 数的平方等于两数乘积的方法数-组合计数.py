"""1577. 数的平方等于两数乘积的方法数"""

from typing import List
from collections import Counter

# 给你两个整数数组 nums1 和 nums2 ，请你返回根据以下规则形成的三元组的数目

# 1 <= nums1.length, nums2.length <= 1000

# 哈希表两数之和思想:用两个map记录平方，然后枚举乘积。


class Solution:
    def numTriplets(self, nums1: List[int], nums2: List[int]) -> int:
        n1, n2 = len(nums1), len(nums2)
        c1, c2 = Counter(), Counter()

        for num in nums1:
            c1[num * num] += 1
        for num in nums2:
            c2[num * num] += 1

        res = 0
        for i in range(n1):
            for j in range(i + 1, n1):
                res += c2[nums1[i] * nums1[j]]
        for i in range(n2):
            for j in range(i + 1, n2):
                res += c1[nums2[i] * nums2[j]]

        return res


print(Solution().numTriplets(nums1=[1, 1], nums2=[1, 1, 1]))
# 输出：9
# 解释：所有三元组都符合题目要求，因为 1^2 = 1 * 1
# 类型 1：(0,0,1), (0,0,2), (0,1,2), (1,0,1), (1,0,2), (1,1,2), nums1[i]^2 = nums2[j] * nums2[k]
# 类型 2：(0,0,1), (1,0,1), (2,0,1), nums2[i]^2 = nums1[j] * nums1[k]
