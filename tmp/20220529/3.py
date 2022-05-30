MOD = int(1e9 + 7)
INF = int(1e20)

from itertools import groupby
from typing import List


class Solution:
    def totalSteps2(self, nums: List[int]) -> int:
        """删掉不是LIS里的元素 求剩下里LIS最长?"""
        n = len(nums)

        isOk = [False] * n
        isOk[0] = True
        preMax = nums[0]
        for i in range(1, n):
            if nums[i] >= preMax:
                isOk[i] = True
                preMax = nums[i]
        bad = [nums[i] for i in range(n) if not isOk[i]]
        print(bad)
        groups = [[char, len(list(group))] for char, group in groupby(isOk)]
        return max((count for flag, count in groups if flag == False), default=0)

    def totalSteps(self, nums: List[int]) -> int:
        """单调栈形成掌控的局面"""

        n = len(nums)
        if n <= 1:
            return 0
        # isOk = [False] * n
        # isOk[0] = True
        # preMax = nums[0]
        # for i in range(n):
        #     if nums[i] >= preMax:
        #         isOk[i] = True
        #         preMax = nums[i]

        res = 0
        cur = 0
        stack = [nums[0]]
        isOk = True
        for i in range(1, n):
            if nums[i] >= stack[-1]:
                stack.append(nums[i])
                cur = 0
            elif i + 1 < n and nums[i] > nums[i + 1]:
                isOk = False
                stack.append(nums[i])
                cur = 0
            else:
                isOk = False
                cur += 1
            res = max(res, cur)

        return res + self.totalSteps(stack) if not isOk else res


print(Solution().totalSteps(nums=[5, 3, 4, 4, 7, 3, 6, 11, 8, 5, 11]))  # 3
print(Solution().totalSteps(nums=[4, 5, 7, 7, 13]))  # 0
print(Solution().totalSteps(nums=[10, 1, 2, 3, 4, 5, 6, 1, 2, 3]))  # 6
print(Solution().totalSteps(nums=[7, 14, 4, 14, 13, 2, 6, 13]))  # 3
print(Solution().totalSteps(nums=[3, 2, 1]))  # 2
