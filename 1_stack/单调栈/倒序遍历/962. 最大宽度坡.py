from typing import List

# 坡是元组 (i, j)，其中  i < j 且 A[i] <= A[j]。这样的坡的宽度为 j - i。
# 单调栈/排序id/线段树&树状数组


class Solution:
    def maxWidthRamp(self, nums: List[int]) -> int:
        stack = []  # 维护左侧下降的点
        res = 0
        for i, num in enumerate(nums):
            if not stack or nums[stack[-1]] > num:
                stack.append(i)

        for j in range(len(nums) - 1, -1, -1):
            while stack and nums[stack[-1]] <= nums[j]:
                res = max(res, j - stack.pop())
        return res

    def maxWidthRamp2(self, nums: List[int]) -> int:
        """O(nlogn) 按照nums[i]递增的顺序遍历id,维护之前看过的i的最小值
        """
        res = 0
        ids = sorted(range(len(nums)), key=nums.__getitem__)
        preMin = int(1e20)
        for i in ids:
            res = max(res, i - preMin)
            preMin = min(preMin, i)
        return res

    def maxWidthRamp3(self, nums: List[int]) -> int:
        """对每个数，找到右侧`最后一个`大于它的数的索引(最大的)

        树状数组/线段树 解法 维护区间最值
        """
        ...


print(Solution().maxWidthRamp([9, 8, 1, 0, 1, 9, 4, 0, 4, 1]))
print(Solution().maxWidthRamp2([9, 8, 1, 0, 1, 9, 4, 0, 4, 1]))
# 输出：7
# 解释：
# 最大宽度的坡为 (i, j) = (2, 9): A[2] = 1 且 A[9] = 1.
