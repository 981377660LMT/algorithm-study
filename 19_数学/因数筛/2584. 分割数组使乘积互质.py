# 给你一个长度为 n 的整数数组 nums ，下标从 0 开始。

# !如果在下标 i 处 分割 数组，其中 0 <= i <= n - 2 ，
# !使前 i + 1 个元素的乘积和剩余元素的乘积互质，则认为该分割 有效 。

# 例如，如果 nums = [2, 3, 3] ，那么在下标 i = 0 处的分割有效，
# 因为 2 和 9 互质，而在下标 i = 1 处的分割无效，因为 6 和 3 不互质。
# 在下标 i = 2 处的分割也无效，因为 i == n - 1 。
# 返回可以有效分割数组的最小下标 i ，如果不存在有效分割，则返回 -1 。

# n == nums.length
# 1 <= n <= 1e4
# 1 <= nums[i] <= 1e6

# !预处理每个素因子的范围+差分求出不合法的区域

from itertools import accumulate
from typing import List
from collections import defaultdict
from prime import EratosthenesSieve


S = EratosthenesSieve(int(1e6) + 10)
F = [S.getPrimeFactors(i) for i in range(int(1e6) + 10)]


class Solution:
    def findValidSplit(self, nums: List[int]) -> int:
        # 起点终点
        mp = defaultdict(list)
        for i, x in enumerate(nums):
            for p in F[x]:
                mp[p].append(i)
        for pos in mp.values():
            pos.sort()

        # 在左边或者右边
        n = len(nums)
        diff = [0] * (n + 1)
        for g in mp.values():
            if len(g) == 1:
                continue
            start, end = g[0], g[-1]
            diff[start] += 1
            diff[end] -= 1

        diff = list(accumulate(diff))
        for i in range(n - 1):
            if diff[i] == 0:
                return i
        return -1


print(Solution().findValidSplit(nums=[4, 7, 8, 15, 3, 5]))
