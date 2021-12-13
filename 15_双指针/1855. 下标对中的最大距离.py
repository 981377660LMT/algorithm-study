from typing import List

# 两个 非递增 的整数数组 nums1​​​​​​ 和 nums2
# 如果该下标对同时满足 i <= j 且 nums1[i] <= nums2[j] ，则称之为 有效 下标对，该下标对的 距离 为 j - i​​ 。​​
# 返回所有 有效 下标对 (i, j) 中的 最大距离 。如果不存在有效下标对，返回 0 。
class Solution:
    def maxDistance(self, nums1: List[int], nums2: List[int]) -> int:
        i, j, res = 0, 0, 0
        while i < len(nums1) and j < len(nums2):
            if nums1[i] > nums2[j]:
                i += 1
            else:
                res = max(res, j - i)
                j += 1
        return res


print(Solution().maxDistance(nums1=[55, 30, 5, 4, 2], nums2=[100, 20, 10, 10, 5]))
# 输出：2
# 解释：有效下标对是 (0,0), (2,2), (2,3), (2,4), (3,3), (3,4) 和 (4,4) 。
# 最大距离是 2 ，对应下标对 (2,4) 。
