from itertools import accumulate
from typing import List
from 每个元素作为最值的影响范围 import getRange


MOD = int(1e9 + 7)

# amazon-oa-min-sum-product
# 求sum(子数组最小值*子数组之和)

# !注意这里不能计算重复元素的影响范围，因此一边开一边闭
# 1. 处理出每个元素作为最小值的影响区间(!注意这里不能计算重复元素的影响范围，因此一边开一边闭)
# 2. 左右边界为 [L,R]，那么他对答案产生的贡献是 v 乘上 [L,R] 内所有`包含 v 的子数组`的元素和的和。


class Solution:
    def totalStrength(self, strength: List[int]) -> int:
        minRange = getRange(strength, isMax=False, isLeftStrict=True, isRightStrict=False)
        p1 = [0] + list(accumulate(strength))
        p2 = [0] + list(accumulate(p1))
        res = 0
        for mid, value in enumerate(strength):
            left, right = minRange[mid]
            sum1 = (p2[right + 2] - p2[mid + 1]) * (mid - left + 1)
            sum2 = (p2[mid + 1] - p2[left]) * (right - mid + 1)
            res = (res + (sum1 - sum2) * value) % MOD
        return res


print(Solution().totalStrength([1, 3, 1, 2]))
print(Solution().totalStrength(strength=[5, 4, 6]))

# print(Solution().totalStrength(strength=[3, 1, 6, 4, 5, 2]))

