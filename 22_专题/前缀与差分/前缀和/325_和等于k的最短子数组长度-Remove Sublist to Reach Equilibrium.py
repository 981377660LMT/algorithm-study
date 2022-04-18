# 移除一段最短的子数组，使得剩下部分严格大于k的个数与严格小于k的个数相等。

# 1. 转化为只含有0，-1，1的数组
# 2. 只需求和为sum(nums)的`最短子数组``

from itertools import accumulate
from typing import List

# todo 细节有问题


class Solution:
    def solve(self, nums: List[int], k: int) -> int:
        n = len(nums)
        nums = [int(n > k) - int(n < k) for n in nums]
        sum_ = sum(nums)

        res = int(1e20)
        last = {0: 0}
        # 最短子数组可以不选，所以要从0开始
        preSum = [0] + list(accumulate(nums))
        for i, cur in enumerate(preSum):
            if cur - sum_ in last:
                res = min(res, i - last[cur - sum_])
            last[cur] = i

        return n - res


print(Solution().solve([5, 9, 7, 8, 2, 4], 5))
print(Solution().solve([1, 2, 3], 4))
print(Solution().solve([0], 0))
