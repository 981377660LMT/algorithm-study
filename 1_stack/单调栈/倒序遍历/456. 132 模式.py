#     1
#     1
#     1    1
#     1    1
# 1   1    1
# 1   1    1
# 是否存在132模式的子序列

# 显然维护一个单增的单调栈即可
from typing import List


class Solution:
    def find132pattern(self, nums: List[int]) -> bool:
        """
        从后往前遍历，维护一个单减的栈
        
        可以找到一对(j,k)使得nums[j]>nums[k]
        记录这个被pop出的nums[k]
        如果之后遇到一个元素比nums[k]小 那么就找到了132模式
        """
        minStack = []
        popped = -int(1e20)
        for num in reversed(nums):
            if num < popped:
                return True
            while minStack and minStack[-1] < num:
                popped = minStack.pop()
            minStack.append(num)
        return False

