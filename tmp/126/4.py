from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个长度为 n 的整数数组 nums 和一个 正 整数 k 。
# 一个整数数组的 能量 定义为和 等于 k 的子序列的数目。
# 请你返回 nums 中所有子序列的 能量和 。
# 由于答案可能很大，请你将它对 109 + 7 取余 后返回。
# 先求有多少种组成k的子集，再看这些子集属于多少个子序列


class Solution:
    def sumOfPower(self, nums: List[int], k: int) -> int:
        dp = defaultdict(int)  # {(sum, count): freq}
        dp[(0, 0)] = 1
        for num in nums:
            for (sum_, count), freq in list(dp.items()):
                if sum_ + num <= k:
                    dp[(sum_ + num, count + 1)] += freq
        counter = defaultdict(int)  # {count: freq}
        for (sum_, count), freq in dp.items():
            if sum_ == k and count > 0:
                counter[count] += freq
        n = len(nums)
        res = 0
        for count, freq in counter.items():
            belong = pow(2, n - count, MOD) * freq
            res += belong
            res %= MOD
        return res % MOD


print(Solution().sumOfPower([1, 2, 3], 3))
