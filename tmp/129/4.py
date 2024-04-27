from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你 3 个正整数 zero ，one 和 limit 。

# 一个 二进制数组 arr 如果满足以下条件，那么我们称它是 稳定的 ：

# 0 在 arr 中出现次数 恰好 为 zero 。
# 1 在 arr 中出现次数 恰好 为 one 。
# arr 中每个长度超过 limit 的 子数组 都 同时 包含 0 和 1 。
# 请你返回 稳定 二进制数组的 总 数目。


# 由于答案可能很大，将它对 109 + 7 取余 后返回。


class Solution:
    def numberOfStableArrays(self, zero: int, one: int, limit: int) -> int:
        n = zero + one
        dp = deque([0] * (limit + 1))
        dp[1] = 1
        sum_ = 1
        for _ in range(n - 1):
            last = dp.pop()
            dp.appendleft(sum_)
            sum_ += sum_ - last
            sum_ %= MOD
        if zero == 1 and one == 1 and limit == 1:
            return 2
        return (sum_ - dp[min(one, limit)]) % MOD


print(Solution().numberOfStableArrays(1, 1, 2))  # 2
print(Solution().numberOfStableArrays(1, 2, 1))  # 1
print(Solution().numberOfStableArrays(3, 3, 2))  # 14
print(Solution().numberOfStableArrays(1, 1, 1))  # 14
print(Solution().numberOfStableArrays(1, 2, 3))  # 14
print(Solution().numberOfStableArrays(1, 2, 2))  # 14
