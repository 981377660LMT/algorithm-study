from typing import List
from collections import Counter
from itertools import accumulate, count

# 给定一个整数数组 A，返回其中元素之和可被 K 整除的（连续、非空）子数组的数目。
# 1 <= A.length <= 30000

# 按照模k分桶
class Solution:
    def subarraysDivByK(self, nums: List[int], k: int) -> int:
        mod = [v % k for v in [0] + list(accumulate(nums))]
        counter = Counter(mod)
        return sum(count * (count - 1) // 2 for _, count in counter.items())


print(Solution().subarraysDivByK([4, 5, 0, -2, -3, 1], 5))
