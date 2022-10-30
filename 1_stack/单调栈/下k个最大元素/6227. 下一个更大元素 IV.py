# https://leetcode.cn/problems/next-greater-element-iv/
# !求每个数右侧第二个严格大于它的数
# 两次单调栈
# !第一个单调栈pop出去的元素放到第二个单调栈里面
# !第二个单调栈再被pop时统计
from typing import List


class Solution:
    def secondGreaterElement(self, nums: List[int]) -> List[int]:
        n = len(nums)
        stack1, stack2 = [], []
        res = [-1] * n
        for i in range(n):
            while stack2 and nums[stack2[-1]] < nums[i]:
                res[stack2.pop()] = nums[i]
            tmp = []
            while stack1 and nums[stack1[-1]] < nums[i]:
                tmp.append(stack1.pop())
            stack1.append(i)
            stack2.extend(tmp[::-1])
        return res


assert Solution().secondGreaterElement(nums=[11, 13, 15, 12, 0, 15, 12, 11, 9]) == [
    15,
    15,
    -1,
    -1,
    12,
    -1,
    -1,
    -1,
    -1,
]
