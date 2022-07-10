from fractions import Fraction
from typing import List, Tuple


class Solution:
    def validSubarraySize(self, nums: List[int], threshold: int) -> int:
        """找到长度为 k 的 nums 子数组，满足数组中 每个 元素都 大于 threshold / k 。
        最小值>threshold/k
        然后看每个元素作为最小值的影响范围
        """
        ranges = getRange(nums, isMax=False, isLeftStrict=False, isRightStrict=False)
        for i, (left, right) in enumerate(ranges):
            # !这里避免浮点数可以用乘法判断
            if nums[i] > Fraction(threshold, (right - left + 1)):
                return right - left + 1
        return -1


def getRange(
    nums: List[int],
    *,
    isMax=False,
    isLeftStrict=True,
    isRightStrict=False,
) -> List[Tuple[int, int]]:
    """
    求每个元素作为最值的影响范围(区间)
    默认为每个数作为左严格右非严格最小值的影响区间 [left,right]

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


print(Solution().validSubarraySize(nums=[1, 3, 4, 3, 1], threshold=6))
print(Solution().validSubarraySize(nums=[5, 6, 7, 8, 9], threshold=24))
print(
    Solution().validSubarraySize(
        [818, 232, 595, 418, 608, 229, 37, 330, 876, 774, 931, 939, 479, 884, 354, 328],
        3790,
    )
)
# -1
