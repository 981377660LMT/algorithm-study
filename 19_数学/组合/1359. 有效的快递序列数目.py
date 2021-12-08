from functools import lru_cache

MOD = int(1e9 + 7)

# 请你统计所有有效的 收件/配送 序列的数目，确保第 i 个物品的配送服务 delivery(i) 总是在其收件服务 pickup(i) 之后。
# 每一项为(2n)!/(2^n)
#  f(n) = (2n^2 - n) * f(n-1);
class Solution:
    def countOrders(self, n: int) -> int:
        res = 1
        for i in range(1, n):
            res *= (2 * i + 1) * (i + 1)
            res %= MOD
        return res


print(Solution().countOrders(n=2))

# 输出：6
# 解释：所有可能的序列包括：
# (P1,P2,D1,D2)，(P1,P2,D2,D1)，(P1,D1,P2,D2)，(P2,P1,D1,D2)，(P2,P1,D2,D1) 和 (P2,D2,P1,D1)。
# (P1,D2,P2,D1) 是一个无效的序列，因为物品 2 的收件服务（P2）不应在物品 2 的配送服务（D2）之后。

