from typing import List, Tuple
from itertools import accumulate


MOD = int(1e9 + 7)

# amazon-oa-min-sum-product
# 求sum(子数组最小值*子数组之和)

# !注意这里不能计算重复元素的影响范围，因此一边开一边闭
# 1. 处理出每个元素作为最小值的影响区间(!注意这里不能计算重复元素的影响范围，因此一边开一边闭)
# 2. 左右边界为 [L,R]，那么他对答案产生的贡献是 v 乘上 [L,R] 内所有`包含 v 的子数组`的元素和的和。
# 推公式
# https://leetcode.cn/problems/sum-of-total-strength-of-wizards/solution/dan-diao-zhan-qian-zhui-he-de-qian-zhui-d9nki/


class Solution:
    def totalStrength(self, strength: List[int]) -> int:
        minRange = getRange(strength, isMax=False, isLeftStrict=True, isRightStrict=False)
        manager = SubArraySumManager(strength)
        res = 0
        for mid, value in enumerate(strength):
            left, right = minRange[mid]
            res += manager.querySubArraySumInclude(left, right, mid) * value
            # 也可以反过来想 [left,right]里所有子数组减去[left,right]里不包含mid的
            # res += (
            #     manager.querySubArraySum(left, right)
            #     - manager.querySubArraySum(left, mid - 1)
            #     - manager.querySubArraySum(mid + 1, right)
            # ) * value
        return res % MOD


def getRange(
    nums: List[int],
    *,
    isMax=False,
    isLeftStrict=False,
    isRightStrict=False,
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


# 巫师的力量和
class SubArraySumManager:
    def __init__(self, nums: List[int]) -> None:
        self.nums = nums
        self.p1 = list(accumulate(nums, initial=0))
        self.p2 = list(accumulate(self.p1, initial=0))
        self.rp1 = (list(accumulate(nums[::-1], initial=0)))[::-1]
        self.rp2 = (list(accumulate(self.rp1[::-1], initial=0)))[::-1]

        self.pi1 = None  # i*nums[i] 前缀和
        self.pi2 = None  # i*i*nums[i] 前缀和

    def querySubArraySum(self, left: int, right: int) -> int:
        """O(1)查询[left,right]闭区间内所有子数组的和

        a[i] 被计算 (i-L+1)*(R-i+1) 次
        [L,R]里所有子数组和为 ∑(L,R) nums[i]*(-i^2+(R+L)*i-LR+R-L+1)
        预处理出 nums[i]、i*nums[i]、i*i*nums[i] 的前缀和即可
        """
        if self.pi1 is None:
            self.pi1 = [0]
            for i, num in enumerate(self.nums):
                self.pi1.append(self.pi1[-1] + num * i)
        if self.pi2 is None:
            self.pi2 = [0]
            for i, num in enumerate(self.nums):
                self.pi2.append(self.pi2[-1] + num * i * i)

        sum1 = self.pi2[left] - self.pi2[right + 1]
        sum2 = (left + right) * (self.pi1[right + 1] - self.pi1[left])
        sum3 = (right + 1 - left - left * right) * (self.p1[right + 1] - self.p1[left])
        return sum1 + sum2 + sum3

    def querySubArraySumStartsAt(self, left: int, right: int) -> int:
        """O(1)查询[left,right]闭区间内所有以nums[left]开头的子数组的和"""
        assert 0 <= left <= right <= len(self.nums) - 1
        sum1 = self.rp1[left] * (right - left + 1)
        sum2 = self.rp2[left + 1] - self.rp2[right + 2]
        return sum1 - sum2

    def querySubArraySumEndsAt(self, left: int, right: int) -> int:
        """O(1)查询[left,right]闭区间内所有以nums[right]结尾的子数组的和

        (p1[right+1]-p1[left])+(p1[right+1]-p1[left+1])+...+(p1[right+1]-p1[right]]) 即
        p1[right+1]*(right-left+1)+p2[right+1]-p2[left]
        """
        assert 0 <= left <= right <= len(self.nums) - 1
        sum1 = self.p1[right + 1] * (right - left + 1)
        sum2 = self.p2[right + 1] - self.p2[left]
        return sum1 - sum2

    def querySubArraySumInclude(self, left: int, right: int, include: int) -> int:
        """O(1)查询[left,right]闭区间内所有包含include下标的子数组的和"""
        assert 0 <= left <= include <= right <= len(self.nums) - 1
        sum1 = (self.p2[right + 2] - self.p2[include + 1]) * (include - left + 1)
        sum2 = (self.p2[include + 1] - self.p2[left]) * (right - include + 1)
        return sum1 - sum2

    def querySubArrayOccurrence(self, left: int, right: int, i: int) -> int:
        """O(1)查询[left,right]闭区间内第i个元素在区间内所有子数组中出现的次数"""
        assert 0 <= left <= right <= len(self.nums) - 1
        return (right - i + 1) * (i - left + 1)


print(Solution().totalStrength([1, 3, 1, 2]))
print(Solution().totalStrength(strength=[5, 4, 6]))
