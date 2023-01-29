#     1
#     1
#     1    1
#     1    1
# 1   1    1
# 1   1    1
# 是否存在132模式的子序列
# 132 模式的子序列 由三个整数 nums[i]、nums[j] 和 nums[k] 组成，
# 并同时满足：i < j < k 和 nums[i] < nums[k] < nums[j] 。
# !n<=2e5

from typing import List
from sortedcontainers import SortedList

INF = int(1e20)


class Solution:
    def find132pattern1(self, nums: List[int]) -> bool:
        """枚举1 单调栈O(n)

        !从后往前遍历，维护一个单减的栈

        可以找到一对(j,k)使得nums[j]>nums[k]
        记录这个被pop出的nums[k]
        如果之后遇到一个元素比nums[k]小 那么就找到了132模式
        """
        minStack = []
        num3 = -INF
        for num1 in reversed(nums):
            if num1 < num3:
                return True
            while minStack and minStack[-1] < num1:
                num3 = minStack.pop()
            minStack.append(num1)
        return False

    def find132pattern3(self, nums: List[int]) -> bool:
        """枚举3(中间的数) O(nlogn)

        对1:维护左侧最小值
        对2:维护右侧有序集合,找到第一个比左侧最小值大的数,检验是否比中间数小
        """
        n = len(nums)
        leftMin = INF
        right = SortedList(nums)
        for i2 in range(n):
            leftMin = min(leftMin, nums[i2 - 1] if i2 > 0 else INF)
            right.remove(nums[i2])
            pos = right.bisect_right(leftMin)
            if pos < len(right) and right[pos] < nums[i2]:
                return True
        return False


print(Solution().find132pattern1([3, 1, 4, 2]))
print(Solution().find132pattern3([-2, 1, -2]))
