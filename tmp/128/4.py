from bisect import bisect_left
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个 正 整数数组 nums 。

# 请你求出 nums 中有多少个子数组，满足子数组中 第一个 和 最后一个 元素都是这个子数组中的 最大 值。


from typing import List, Tuple


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


class Solution:
    def numberOfSubarrays(self, nums: List[int]) -> int:
        ranges = getRange(nums, isMax=True, isLeftStrict=False, isRightStrict=True)
        mp = defaultdict(list)
        for i, v in enumerate(nums):
            mp[v].append(i)

        def query(start: int, end: int, value: int) -> int:
            return bisect_left(mp[value], end) - bisect_left(mp[value], start)

        res = 0
        for i, (left, right) in enumerate(ranges):
            cur = nums[i]
            res += query(left, right + 1, cur)
            # print(cur, left, right, query(left, right + 1, cur), res)
        return res


print(Solution().numberOfSubarrays([3, 2, 3]))
# nums = [1,4,3,3,2]
print(Solution().numberOfSubarrays([1, 4, 3, 3, 2]))
