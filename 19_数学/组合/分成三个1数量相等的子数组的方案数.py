# 分成三个1数量相等的子数组的方案数
from collections import Counter
from itertools import accumulate
from math import comb


class Solution:
    def solve(self, s):
        nums = list(map(int, s))
        preSum = [0] + list(accumulate(nums))
        sum_ = preSum[-1]
        n = len(preSum)
        MOD = int(1e9 + 7)
        counter = Counter(preSum)  # 注意这里

        if sum_ % 3 != 0:
            return 0

        # 注意组合数相等的特判
        if sum_ == 0:
            return comb(n - 2, 2) % MOD

        return counter[sum_ - sum_ // 3] * counter[sum_ - 2 * sum_ // 3] % MOD


print(Solution().solve(s="11001111"))
# We can have

# "11" + "0011" + "11"
# "110" + "011" + "11"
# "1100" + "11" + "11"
