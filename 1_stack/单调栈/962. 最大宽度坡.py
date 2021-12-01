from typing import List

# 坡是元组 (i, j)，其中  i < j 且 A[i] <= A[j]。这样的坡的宽度为 j - i。
# 一眼单调栈
class Solution:
    def maxWidthRamp(self, nums: List[int]) -> int:
        stack = []
        res = 0
        for i, num in enumerate(nums):
            if not stack or nums[stack[-1]] > num:
                stack.append(i)
        print(stack)
        for j in range(len(nums) - 1, -1, -1):
            while stack and nums[stack[-1]] <= nums[j]:
                res = max(res, j - stack.pop())
        return res


print(Solution().maxWidthRamp([9, 8, 1, 0, 1, 9, 4, 0, 4, 1]))
# 输出：7
# 解释：
# 最大宽度的坡为 (i, j) = (2, 9): A[2] = 1 且 A[9] = 1.
