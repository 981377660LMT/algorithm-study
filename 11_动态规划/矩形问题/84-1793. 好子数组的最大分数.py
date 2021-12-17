from typing import List

# 1 <= nums.length <= 105
# 直方图中最大矩形面积变形题

# 一个 好 子数组的两个端点下标需要满足 i <= k <= j 。
# 请你返回 好 子数组的最大可能 分数 。

# 双指针
# 84题，多了一个判断更新res的条件
class Solution:
    def maximumScore(self, nums: List[int], k: int) -> int:
        nums = [-1] + nums + [-1]
        res = 0
        stack = []
        for i, num in enumerate(nums):
            while stack and nums[stack[-1]] > num:
                height = nums[stack.pop()]

                left = stack[-1] + 1
                right = i - 1
                # 注意还原下标到原数组,因为多了两个哨兵
                if left <= k + 1 <= right:
                    res = max(res, height * (right - left + 1))

            stack.append(i)
        return res


print(Solution().maximumScore(nums=[1, 4, 3, 7, 4, 5], k=3))
# 输出：15
# 解释：最优子数组的左右端点下标是 (1, 5) ，分数为 min(4,3,7,4,5) * (5-1+1) = 3 * 5 = 15 。
