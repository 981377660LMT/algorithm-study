from typing import List
from functools import lru_cache
from sys import setrecursionlimit

# 你有两个 有序 且数组内元素互不相同的数组 nums1 和 nums2 。
# 如果你遇到了 nums1 和 nums2 中都存在的值，那么你可以切换路径到另一个数组对应数字处继续遍历
# 得分定义为合法路径中不同数字的和。
# 请你返回所有可能合法路径中的最大得分。
# 由于答案可能很大，请你将它对 10^9 + 7 取余后返回。
# 1 <= nums1.length <= 10^5
# 1 <= nums2.length <= 10^5


# 总结：
# https://leetcode-cn.com/problems/get-the-maximum-score/solution/pythonshuang-zhi-zhen-yi-ci-bian-li-xiang-xi-tu-ji/


MOD = int(1e9 + 7)


class Solution:
    def maxSum(self, nums1: List[int], nums2: List[int]) -> int:
        m, n = len(nums1), len(nums2)
        sum1, sum2 = 0, 0
        i, j = 0, 0

        while i < m and j < n:
            # 小的一方先走
            if nums1[i] < nums2[j]:
                sum1 += nums1[i]
                i += 1
            elif nums2[j] < nums1[i]:
                sum2 += nums2[j]
                j += 1
            else:
                # 每遇到岔口，更新两个sum为最大值，同时取余
                sum1 = sum2 = (max(sum1, sum2) + nums1[i]) % MOD
                i += 1
                j += 1

        sum1, sum2 = sum1 + sum(nums1[i:]), sum2 + sum(nums2[j:])
        return max(sum1, sum2) % MOD


# print(Solution().maxSum(nums1=[2, 4, 5, 8, 10], nums2=[4, 6, 8, 9]))
# 输出：30
# 解释：合法路径包括：
# [2,4,5,8,10], [2,4,5,8,9], [2,4,6,8,9], [2,4,6,8,10],（从 nums1 开始遍历）
# [4,6,8,9], [4,5,8,10], [4,5,8,9], [4,6,8,10]  （从 nums2 开始遍历）
# 最大得分为上图中的绿色路径 [2,4,6,8,10] 。
print(Solution().maxSum(nums1=[1, 3, 5, 7, 9], nums2=[3, 5, 100]))

