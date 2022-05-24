from typing import List, Tuple
from SubArraySumManager import SubArraySumManager

MOD = int(1e9 + 7)

# amazon-oa-min-sum-product
# 求sum(子数组最小值*子数组之和)

# !注意这里不能计算重复元素的影响范围，因此一边开一边闭
# 1. 处理出每个元素作为最小值的影响区间(!注意这里不能计算重复元素的影响范围，因此一边开一边闭)
# 2. 左右边界为 [L,R]，那么他对答案产生的贡献是 v 乘上 [L,R] 内所有`包含 v 的子数组`的元素和的和。
# 推公式
# https://leetcode.cn/problems/sum-of-total-strength-of-wizards/solution/dan-diao-zhan-qian-zhui-he-de-qian-zhui-d9nki/


MOD = int(1e9 + 7)


class Solution:
    def totalStrength(self, strength: List[int]) -> int:
        minRange = getRange(strength, isMax=False, isLeftStrict=True, isRightStrict=False)
        manager = SubArraySumManager(strength)
        res = 0
        for mid, value in enumerate(strength):
            left, right = minRange[mid]
            # res += manager.querySubArraySumInclude(left, right, mid) * value
            # 也可以反过来想 [left,right]里所有子数组减去[left,right]里不包含mid的
            res += (
                manager.querySubArraySum(left, right)
                - manager.querySubArraySum(left, mid - 1)
                - manager.querySubArraySum(mid + 1, right)
            ) * value
        return res % MOD


def getRange(
    nums: List[int], *, isMax=False, isLeftStrict=False, isRightStrict=False,
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


print(Solution().totalStrength([1, 3, 1, 2]))
print(Solution().totalStrength(strength=[5, 4, 6]))
