# 每次 操作 会从中选择一个元素并 将该元素的值减少 1。
# 返回将数组 nums 转换为锯齿数组所需的最小操作次数。
# !考察奇数位/偶数位

from typing import List


# !Two options, either make A[even] smaller or make A[odd] smaller.


class Solution:
    def movesToMakeZigzag(self, nums: List[int]) -> int:
        n = len(nums)
        cost = [0, 0]
        for start in (0, 1):
            for i in range(start, n, 2):
                diff = 0
                # 比较两边，如果大了，就要减
                if i > 0:
                    diff = max(diff, nums[i] - nums[i - 1] + 1)
                if i + 1 < n:
                    diff = max(diff, nums[i] - nums[i + 1] + 1)
                cost[start] += diff
        return min(cost)


print(Solution().movesToMakeZigzag(nums=[9, 6, 1, 6, 2]))
# 输出：4
