"""
# 6264. 让数组不相等的最小总代价
https://leetcode.cn/problems/minimum-total-cost-to-make-arrays-unequal/
每次操作中，你可以选择交换 nums1 中任意两个下标处的值。操作的 开销 为两个下标的 和 。
你的目标是对于所有的 0 <= i <= n - 1 ，都满足 nums1[i] != nums2[i] ，
你可以进行 任意次 操作，请你返回达到这个目标的 最小 总代价。
请你返回让 nums1 和 nums2 满足上述条件的 最小总代价 ，如果无法达成目标，返回 -1 。

偶数对相同:可以直接交换
奇数对相同:可以将0作为中转站,变为偶数情况
"""

# !贪心的从小到大去dilute不对 还要跳过和当前maxCount相同的下标
# 还要跳过nums1[i]==maxValnum或者nums2[i]==maxVal的下标。

from typing import List
from collections import defaultdict, Counter


class Solution:
    def minimumTotalCost(self, nums1: List[int], nums2: List[int]) -> int:
        n = len(nums1)
        counter = Counter(nums1 + nums2)
        if any(v > n for v in counter.values()):
            return -1

        sameCounter = defaultdict(int)
        maxCount = 0
        count = 0
        visited = [False] * n
        res = []
        for i, (num1, num2) in enumerate(zip(nums1, nums2)):
            if num1 == num2:
                sameCounter[num1] += 2
                count += 2
                maxCount = max(maxCount, sameCounter[num1])
                visited[i] = True
                res.append(i)

        if maxCount <= count // 2:
            return sum(res)

        for i in range(n):
            if visited[i] or sameCounter[nums1[i]] == maxCount or sameCounter[nums2[i]] == maxCount:
                continue
            for num in (nums1[i], nums2[i]):
                sameCounter[num] += 1
                maxCount = max(maxCount, sameCounter[num])
                count += 1
            res.append(i)
            if maxCount <= count // 2:
                return sum(res)


print(Solution().minimumTotalCost([1, 2, 3, 4, 5], [1, 2, 3, 4, 5]))
# [1,5,3,5,5]
# [1,2,3,4,5]
print(Solution().minimumTotalCost([1, 5, 3, 5, 5], [1, 2, 3, 4, 5]))
print(Solution().minimumTotalCost([1, 2, 3, 4], [2, 2, 2, 2]))
# [2,2,2,3,4,5]
# [5,2,2,3,1,4]
print(Solution().minimumTotalCost([2, 2, 2, 3, 4, 5], [5, 2, 2, 3, 1, 4]))
# [8,8,1,2,2,3,2,6,1,10] [8,4,6,3,7,1,5,4,9,5]
print(Solution().minimumTotalCost([8, 8, 1, 2, 2, 3, 2, 6, 1, 10], [8, 4, 6, 3, 7, 1, 5, 4, 9, 5]))
# !贪心的从小到大去dilute不对 还要跳过和当前maxCount相同的下标
# 还要跳过nums1[i]==maxValnum或者nums2[i]==maxVal的下标。
