from typing import List

# 记录两个数组前缀和差值第一次出现的索引位置即可
class Solution:
    def widestPairOfIndices(self, nums1: List[int], nums2: List[int]) -> int:
        res = 0
        n = len(nums1)
        pre1, pre2 = 0, 0

        diff = {0: -1}
        for i in range(n):
            pre1 += 1 if nums1[i] else 0
            pre2 += 1 if nums2[i] else 0
            if pre2 - pre1 in diff:
                res = max(res, i - diff[pre2 - pre1])
            else:
                diff[pre2 - pre1] = i

        return res


print(Solution().widestPairOfIndices(nums1=[1, 1, 0, 1], nums2=[0, 1, 1, 0]))

# Output: 3
# Explanation:
# If i = 1 and j = 3:
# nums1[1] + nums1[2] + nums1[3] = 1 + 0 + 1 = 2.
# nums2[1] + nums2[2] + nums2[3] = 1 + 1 + 0 = 2.
# The distance between i and j is j - i + 1 = 3 - 1 + 1 = 3.
