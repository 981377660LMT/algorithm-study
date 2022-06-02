from typing import List, Tuple


# 给你一个整数数组 nums 。nums 中，子数组的 范围 是子数组中最大元素和最小元素的差值。
# 返回 nums 中 所有 子数组范围的 和 。


# O(n^2) dp
class Solution1:
    def subArrayRanges(self, nums: List[int]) -> int:
        res = 0
        n = len(nums)
        for i in range(n):
            minVal = nums[i]
            maxVal = nums[i]
            for j in range(i + 1, n):
                minVal = min(minVal, nums[j])
                maxVal = max(maxVal, nums[j])
                res += maxVal - minVal
        return res


def getRange(
    nums: List[int], *, isMax=False, isLeftStrict=True, isRightStrict=False,
) -> List[Tuple[int, int]]:
    """
    求每个元素作为最值的影响范围(区间)
    默认为每个数作为非严格最小值的影响区间 [left,right]

    有时候为了避免重复计算，我们可以考虑左侧`严格小于`当前元素的最近元素位置，
    以及右侧`小于等于`当前元素的最近元素位置。
    """

    def compareLeft(stackValue: int, curValue: int) -> bool:
        if isLeftStrict and isMax:
            return stackValue <= curValue
        elif isLeftStrict and not isMax:
            return stackValue >= curValue
        elif not isLeftStrict and isMax:
            return stackValue < curValue
        else:
            return stackValue > curValue

    def compareRight(stackValue: int, curValue: int) -> bool:
        if isRightStrict and isMax:
            return stackValue <= curValue
        elif isRightStrict and not isMax:
            return stackValue >= curValue
        elif not isRightStrict and isMax:
            return stackValue < curValue
        else:
            return stackValue > curValue

    n = len(nums)
    leftMost = [0] * n
    rightMost = [n - 1] * n

    stack = []
    for i in range(n):
        while stack and compareRight(nums[stack[-1]], nums[i]):
            rightMost[stack.pop()] = i - 1
        stack.append(i)

    stack = []
    for i in range(n - 1, -1, -1):
        while stack and compareLeft(nums[stack[-1]], nums[i]):
            leftMost[stack.pop()] = i + 1
        stack.append(i)

    return list(zip(leftMost, rightMost))


# 单调栈
class Solution:
    def subArrayRanges(self, nums: List[int]) -> int:
        """每个元素作为最值影响了多少个子数组"""
        minRange = getRange(nums, isMax=False, isLeftStrict=True, isRightStrict=False)
        maxRange = getRange(nums, isMax=True, isLeftStrict=True, isRightStrict=False)
        res = 0
        for i in range(len(nums)):
            left1, right1 = minRange[i]
            res -= (i - left1 + 1) * (right1 - i + 1) * nums[i]
            left2, right2 = maxRange[i]
            res += (i - left2 + 1) * (right2 - i + 1) * nums[i]
        return res


print(Solution().subArrayRanges(nums=[1, 3, 3]))
print(Solution().subArrayRanges(nums=[4, -2, -3, 4, 1]))

