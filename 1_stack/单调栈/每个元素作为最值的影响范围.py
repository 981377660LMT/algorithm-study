# 计算每个元素作为最值的影响范围
from typing import List, Tuple


def getRange(nums: List[int], isStrict=False, asMax=False) -> List[Tuple[int, int]]:
    """
    求每个元素作为最值的影响范围(区间)
    默认为每个数作为非严格最小值的影响区间 [left,right]
    """

    def compare(stackValue: int, curValue: int) -> bool:
        if isStrict and asMax:
            return stackValue <= curValue
        elif isStrict and not asMax:
            return stackValue >= curValue
        elif not isStrict and asMax:
            return stackValue < curValue
        else:
            return stackValue > curValue

    n = len(nums)
    leftMost = [0] * n
    rightMost = [n - 1] * n

    stack = []
    for i in range(n):
        while stack and compare(nums[stack[-1]], nums[i]):
            rightMost[stack.pop()] = i - 1
        stack.append(i)

    stack = []
    for i in range(n - 1, -1, -1):
        while stack and compare(nums[stack[-1]], nums[i]):
            leftMost[stack.pop()] = i + 1
        stack.append(i)

    return list(zip(leftMost, rightMost))


if __name__ == '__main__':
    print(getRange([0, 10, 20, 20, 50, 10]))

