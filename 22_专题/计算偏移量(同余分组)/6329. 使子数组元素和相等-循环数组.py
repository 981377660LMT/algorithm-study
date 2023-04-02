# 6329. 使子数组元素和相等-循环数组周期分组

# 给定一个无限循环数组
# !要求所有长为k的子数组元素和相等
# 求+1-1的最小操作次数


# 分组问题:
# 1. 并查集肯定可以
# 2. 但是这里的分组是有规律的,
# !   一个循环数组如果既有周期 n 也有周期 k, 那么必然有周期 gcd(n,k)


from collections import defaultdict
from math import gcd
from typing import List


class Solution:
    def makeSubKSumEqual(self, arr: List[int], k: int) -> int:
        n = len(arr)
        gcd_ = gcd(n, k)
        mp = defaultdict(list)
        for i in range(n):
            mp[i % gcd_].append(arr[i])
        res = 0
        for g in mp.values():
            g.sort()
            mid = g[len(g) // 2]
            res += sum(abs(x - mid) for x in g)
        return res


assert Solution().makeSubKSumEqual([1, 4, 1, 3], 2) == 1
