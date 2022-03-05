from typing import List


class Solution:
    def subArrayRanges(self, nums: List[int]) -> int:
        n = len(nums)
        minLeft, maxLeft = [0] * n, [0] * n
        minStack, maxStack = [], []
        for i, num in enumerate(nums):
            while minStack and nums[minStack[-1]] > num:
                minStack.pop()
            minLeft[i] = minStack[-1] if minStack else -1
            minStack.append(i)

            # 如果 nums[maxStack[-1]] == num, 那么根据定义，
            # nums[maxStack[-1]] 逻辑上小于 num，因为 maxStack[-1] < i
            while maxStack and nums[maxStack[-1]] <= num:
                maxStack.pop()
            maxLeft[i] = maxStack[-1] if maxStack else -1
            maxStack.append(i)

        minRight, maxRight = [0] * n, [0] * n
        minStack, maxStack = [], []
        for i in range(n - 1, -1, -1):
            num = nums[i]
            # 如果 nums[minStack[-1]] == num, 那么根据定义，
            # nums[minStack[-1]] 逻辑上大于 num，因为 minStack[-1] > i
            while minStack and nums[minStack[-1]] >= num:
                minStack.pop()
            minRight[i] = minStack[-1] if minStack else n
            minStack.append(i)

            while maxStack and nums[maxStack[-1]] < num:
                maxStack.pop()
            maxRight[i] = maxStack[-1] if maxStack else n
            maxStack.append(i)

        sumMax, sumMin = 0, 0
        for i, num in enumerate(nums):
            sumMax += (maxRight[i] - i) * (i - maxLeft[i]) * num
            sumMin += (minRight[i] - i) * (i - minLeft[i]) * num
        return sumMax - sumMin
