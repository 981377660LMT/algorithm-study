from itertools import accumulate
from typing import List, Tuple
from collections import defaultdict

MOD = int(1e9 + 7)
INF = int(1e20)

# amazon-oa-min-sum-product
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


class Solution:
    def totalStrength(self, strength: List[int]) -> int:
        # preSum = [0] + list(accumulate(strength))
        ranges = getRange(strength, isStrict=False, asMax=False)
        res = 0
        n = len(strength)
        for mid in range(n):
            left, right = ranges[mid]
            for j in range(left, right + 1):
                if j <= mid:
                    count = (j - left + 1) * (mid - left + 1)
                else:
                    count = (right - mid + 1) * (j - right + 1)
                res += count * strength[j] * strength[mid]
                res %= MOD
        return res


print(Solution().totalStrength([1, 3, 1, 2]))
print(Solution().totalStrength(strength=[5, 4, 6]))
# print(Solution().totalStrength(strength=[3, 1, 6, 4, 5, 2]))

