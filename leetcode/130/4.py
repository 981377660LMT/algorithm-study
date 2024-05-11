from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 一个整数 x 的 强数组 指的是满足和为 x 的二的幂的最短有序数组。比方说，11 的强数组为 [1, 2, 8] 。

# 我们将每一个正整数 i （即1，2，3等等）的 强数组 连接得到数组 big_nums ，big_nums 开始部分为 [1, 2, 1, 2, 4, 1, 4, 2, 4, 1, 2, 4, 8, ...] 。

# 给你一个二维整数数组 queries ，其中 queries[i] = [fromi, toi, modi] ，你需要计算 (big_nums[fromi] * big_nums[fromi + 1] * ... * big_nums[toi]) % modi 。


# 请你返回一个整数数组 answer ，其中 answer[i] 是第 i 个查询的答案。


def cal(upper: int, k: int) -> int:
    """[0,upper]中二进制第k(k>=0)位为1的数的个数.
    即满足 `num & (1 << k) > 0` 的数的个数
    """
    bit = upper.bit_length()
    if k >= bit:
        return 0

    @lru_cache(None)
    def dfs(pos: int, hasLeadingZero: bool, isLimit: bool) -> int:
        """当前在第pos位,hasLeadingZero表示有前导0，isLimit表示是否贴合上界"""
        if pos == -1:
            return not hasLeadingZero
        if pos == k:  # 这一位必须填1
            if isLimit and nums[pos] == 0:
                return 0
            return dfs(pos - 1, False, isLimit)
        res = 0
        up = nums[pos] if isLimit else 1
        for cur in range(up + 1):
            if hasLeadingZero and cur == 0:
                res += dfs(pos - 1, True, (isLimit and cur == up))
            else:
                res += dfs(pos - 1, False, (isLimit and cur == up))
        return res

    nums = list(map(int, bin(upper)[2:]))[::-1]
    res = dfs(len(nums) - 1, True, True)
    dfs.cache_clear()
    return res


BIT_COUNT_SUM = [1, 3, 8, 20, 48, 112]  # BIT_COUNT_SUM[i]表示[2^i,2^(i+1))中1的个数
for _ in range(60):
    i = len(BIT_COUNT_SUM) - 1
    BIT_COUNT_SUM.append(BIT_COUNT_SUM[-1] * 2 + 2**i)


def locateKthBit(k: int) -> Tuple[int, int]:
    """
    给定整数k(k>=0), 返回二进制数中第k个bit所在的数以及k在该数中是第几个1.
    """

    bit = k.bit_length()

    left, right = 0, k
    bitCount = 0
    while left <= right:
        mid = (left + right) // 2
        # check
        sum_ = sum(cal(mid, i) for i in range(bit))
        if sum_ >= k:
            right = mid - 1
            bitCount = sum_
        else:
            left = mid + 1
    print(sum(cal(left, i) for i in range(bit)))
    return left, bitCount - k


print(locateKthBit(6))


class Solution:
    def findProductsOfElements(self, queries: List[List[int]]) -> List[int]:
        res = []
        for a, b, mod in queries:
            cur = 1
            maxBit = b.bit_length()
            for bit in range(maxBit + 1):
                count = cal(b, bit) - cal(a - 1, bit)
                if count == 0:
                    continue
                cur *= pow(1 << bit, count, mod)
                cur %= mod
            print(cur)
            res.append(cur % mod)
        return res


# queries = [[1,3,7]]
# queries = [[2,5,3],[7,7,4]]
print(Solution().findProductsOfElements([[2, 5, 3], [7, 7, 4]]))
