"""
比例子数组的个数
"""

# !公式变形+前缀和
# oneI-oneJ / zeroI-zeroJ = num1 / num2
# num2*(oneI-oneJ) = num1*(zeroI-zeroJ)
# !num2*oneI - num1*zeroI = num2*oneJ - num1*zeroJ (统计i、j对数)
from collections import defaultdict


class Solution:
    def fixedRatio(self, s: str, num1: int, num2: int) -> int:
        """求子数组个数,1的个数/0的个数 = num1/num2"""
        nums = list(map(int, s))
        res, counter, preSum = 0, [0, 0], defaultdict(int, {0: 1})
        for num in nums:
            counter[num] += 1
            cur = num2 * counter[0] - num1 * counter[1]
            res += preSum[cur]
            preSum[cur] += 1
        return res


print(Solution().fixedRatio(s="0110011", num1=1, num2=2))
