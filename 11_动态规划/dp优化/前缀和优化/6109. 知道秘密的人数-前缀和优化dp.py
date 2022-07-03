# 2 <= n <= 1000
# 1 <= delay < forget <= n
from collections import defaultdict

MOD = int(1e9 + 7)


class Solution:
    def peopleAwareOfSecret(self, n: int, delay: int, forget: int) -> int:
        """哈希表模拟过程 O(n)"""
        toBorn = defaultdict(int, {1 + delay: 1})  # 传播者人数
        toDie = defaultdict(int, {1 + forget: 1})  # 传播者死去的人数
        alive = 1  # 当前人数
        spread = 0  # 传播者人数
        for day in range(1, n + 1):
            spread += toBorn[day] - toDie[day]
            newBorn = spread
            toBorn[day + delay] += newBorn
            toDie[day + forget] += newBorn
            alive += newBorn - toDie[day]
            alive %= MOD
        return alive

    def peopleAwareOfSecret2(self, n: int, delay: int, forget: int) -> int:
        """前缀和dp"""
        toBorn = defaultdict(int, {1 + delay: 1})  # 传播者人数
        toDie = defaultdict(int, {1 + forget: 1})  # 传播者死去的人数
        alive = 1  # 当前人数
        spread = 0  # 传播者人数
        for day in range(1, n + 1):
            spread += toBorn[day] - toDie[day]
            newBorn = spread
            toBorn[day + delay] += newBorn
            toDie[day + forget] += newBorn
            alive += newBorn - toDie[day]
            alive %= MOD
        return alive


print(Solution().peopleAwareOfSecret(n=6, delay=2, forget=4))
# class Solution:
#     def peopleAwareOfSecret(self, n: int, delay: int, forget: int) -> int:
#         mod = 10**9+7
#         dp = [0]*forget
#         dp[-1] = 1
#         for _ in range(1,n):
#             cur = sum(dp[-delay:-forget:-1]) % mod
#             dp.append(cur)
#         return sum(dp[-forget:]) % mod
