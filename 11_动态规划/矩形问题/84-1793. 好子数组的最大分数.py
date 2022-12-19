from typing import List

# 1 <= nums.length <= 1e5
# 直方图中最大矩形面积变形题

# 一个子数组 (i, j) 的 分数 定义为
# !min(nums[i], nums[i+1], ..., nums[j]) * (j - i + 1) 。
# !一个 好 子数组的两个端点下标需要满足 i <= k <= j 。
# 请你返回 好 子数组的最大可能 分数 。


class Solution:
    def maximumScore(self, nums: List[int], k: int) -> int:
        n = len(nums)
        leftMost, rightMost = [0] * n, [n - 1] * n

        stack = []
        for i in range(n):
            while stack and nums[stack[-1]] > nums[i]:
                rightMost[stack.pop()] = i - 1
            stack.append(i)

        stack = []
        for i in range(n - 1, -1, -1):
            while stack and nums[stack[-1]] > nums[i]:
                leftMost[stack.pop()] = i + 1
            stack.append(i)

        res = 0
        for i in range(n):
            left, right = leftMost[i], rightMost[i]
            if left <= k <= right:
                res = max(res, nums[i] * (right - left + 1))
        return res


assert Solution().maximumScore(nums=[1, 4, 3, 7, 4, 5], k=3) == 15
# 输出：15
# 解释：最优子数组的左右端点下标是 (1, 5) ，分数为 min(4,3,7,4,5) * (5-1+1) = 3 * 5 = 15 。
