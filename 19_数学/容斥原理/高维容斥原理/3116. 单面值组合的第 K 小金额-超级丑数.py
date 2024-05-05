# 3116. 单面值组合的第 K 小金额-超级丑数
#
# 给你一个整数数组 coins 表示不同面额的硬币，另给你一个整数 k 。
# 你有无限量的每种面额的硬币。但是，你 不能 组合使用不同面额的硬币。
# 返回使用这些硬币能制造的 第 kth 小 金额。

from typing import List
from 能被primes中的至少一个数整除的数的个数 import countMultiple


class Solution:
    def findKthSmallest(self, coins: List[int], k: int) -> int:
        def check(mid: int) -> bool:
            return countMultiple(mid, coins, unique=True) >= k

        left, right = 1, k * min(coins)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left
