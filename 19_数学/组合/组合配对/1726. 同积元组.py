# freq table
# 请你返回满足 a * b = c * d 的元组 (a, b, c, d) 的数量。其中 a、b、c 和 d 都是 nums 中的元素，且 a != b != c != d 。

# 如果 a*b=c*d 那么可以凑出`八对`
# 2x2x2

from collections import Counter
from typing import List
from math import comb


class Solution:
    def tupleSameProduct(self, nums: List[int]) -> int:
        res = 0
        C = Counter(nums[i] * nums[j] for i in range(len(nums)) for j in range(i + 1, len(nums)))

        for count in C.values():
            res += comb(count, 2) * 8
        return res


print(Solution().tupleSameProduct([2, 3, 4, 6]))
# 解释：存在 8 个满足题意的元组：
# (2,6,3,4) , (2,6,4,3) , (6,2,3,4) , (6,2,4,3)
# (3,4,2,6) , (4,3,2,6) , (3,4,6,2) , (4,3,6,2)

