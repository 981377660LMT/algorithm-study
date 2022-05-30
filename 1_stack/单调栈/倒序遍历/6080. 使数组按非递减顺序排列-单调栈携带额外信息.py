from typing import List

# 在一步操作中，移除所有满足 nums[i - 1] > nums[i] 的 nums[i] ，其中 0 < i < nums.length 。
# 重复执行步骤，直到 nums 变为 `非递减` 数组，返回所需执行的操作数。
# 1 <= nums.length <= 1e5

# 1、从后往前遍历，将元素入栈，遇到更大的数时取出
# 2、两段之间合并时，取前一段+1和后一段的最大值，这两段并行删除取最大值


# https://leetcode.cn/problems/steps-to-make-array-non-decreasing/solution/by-qqwqert007-6f65/


class Solution:
    def totalSteps(self, nums: List[int]) -> int:
        """
        反向遍历 维护一个单减的栈
        且为每个栈中的元素维护一个数值，表示共向后删除了几个数
        """
        n = len(nums)
        stack = []
        res = 0
        for i in range(n - 1, -1, -1):
            curRemove = 0  # 这个数删了多少个数
            while stack and stack[-1][0] < nums[i]:
                # 删除是可接力的
                _, preRemove = stack.pop()
                # 两个峰负责的删除可以并行进行，总时间取决于两者中较大者
                curRemove = max(curRemove + 1, preRemove)
            stack.append((nums[i], curRemove))
            res = max(res, curRemove)

        return res


print(Solution().totalSteps([7, 14, 4, 14, 13, 2, 6, 13]))  # 3

