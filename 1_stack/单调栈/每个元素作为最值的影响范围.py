"""
计算每个元素作为最值的影响范围
calculate the range of influence for each element as the maximum/minimum value

结论:
# !1.以元素nums[i]为最值的影响范围[left,right],则包含nums[i]的子数组个数为(right-i+1)*(i-left+1)
"""

from typing import Callable, List, Tuple


def getRange(
    nums: List[int],
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

    def createCompareLeft(isLeftStrict: bool, isMax: bool) -> Callable[[int, int], bool]:
        if isLeftStrict and isMax:
            return lambda stackValue, curValue: stackValue <= curValue
        elif isLeftStrict and not isMax:
            return lambda stackValue, curValue: stackValue >= curValue
        elif not isLeftStrict and isMax:
            return lambda stackValue, curValue: stackValue < curValue
        else:
            return lambda stackValue, curValue: stackValue > curValue

    def createCompareRight(isRightStrict: bool, isMax: bool) -> Callable[[int, int], bool]:
        if isRightStrict and isMax:
            return lambda stackValue, curValue: stackValue <= curValue
        elif isRightStrict and not isMax:
            return lambda stackValue, curValue: stackValue >= curValue
        elif not isRightStrict and isMax:
            return lambda stackValue, curValue: stackValue < curValue
        else:
            return lambda stackValue, curValue: stackValue > curValue

    compareLeft = createCompareLeft(isLeftStrict, isMax)
    compareRight = createCompareRight(isRightStrict, isMax)

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


if __name__ == "__main__":
    print(getRange([0, 10, 20, 20, 50, 10]))
