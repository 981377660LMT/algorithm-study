from math import comb
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 车尔尼有一个数组 nums ，它只包含 正 整数，所有正整数的数位长度都 相同 。

# 两个整数的 数位不同 指的是两个整数 相同 位置上不同数字的数目。


# 请车尔尼返回 nums 中 所有 整数对里，数位不同之和。
class Solution:
    def sumDigitDifferences(self, nums: List[int]) -> int:
        counters = [[0] * 10 for _ in range(len(str(nums[0])))]
        for num in nums:
            for i, c in enumerate(str(num)):
                counters[i][int(c)] += 1
        res = 0
        for counter in counters:
            res += comb(len(nums), 2)
            for c in counter:
                res -= comb(c, 2)
        return res


# nums = [13,23,12]
print(Solution().sumDigitDifferences([13, 23, 12]))


def checkWithBruce(nums: List[int]) -> int:
    res = 0
    for a in nums:
        for b in nums:
            res += sum(1 for x, y in zip(str(a), str(b)) if x != y)
    return res // 2


for _ in range(100):
    import random

    nums = [random.randint(100, 999) for _ in range(3)]
    assert Solution().sumDigitDifferences(nums) == checkWithBruce(nums)

print(checkWithBruce([13, 23, 12]))
