# 100047. 统计树中的合法路径数目
#
# https://leetcode.cn/problems/count-valid-paths-in-a-tree/description/
# !求数中的合法路径数目(路径长度>=2)，路径上的质数个数恰好为1.

from typing import List, Tuple
from Rerooting import Rerooting
from collections import Counter


class EratosthenesSieve:
    """埃氏筛"""

    __slots__ = "_minPrime"  # 每个数的最小质因数

    def __init__(self, maxN: int):
        """预处理 O(nloglogn)"""
        minPrime = list(range(maxN + 1))
        upper = int(maxN**0.5) + 1
        for i in range(2, upper):
            if minPrime[i] < i:
                continue
            for j in range(i * i, maxN + 1, i):
                if minPrime[j] == j:
                    minPrime[j] = i
        self._minPrime = minPrime

    def isPrime(self, n: int) -> bool:
        if n < 2:
            return False
        return self._minPrime[n] == n

    def getPrimeFactors(self, n: int) -> "Counter[int]":
        """求n的质因数分解 O(logn)"""
        res, f = Counter(), self._minPrime
        while n > 1:
            m = f[n]
            res[m] += 1
            n //= m
        return res

    def getPrimes(self) -> List[int]:
        return [x for i, x in enumerate(self._minPrime) if i >= 2 and i == x]


P = EratosthenesSieve(int(1e5) + 10)


class Solution:
    def countPaths(self, n: int, edges: List[List[int]]) -> int:
        E = Tuple[int, int]

        def e(root: int) -> E:
            return (0, 0)  # 幺元 (恰好0个质数的路径数, 恰好1个质数的路径数)

        def op(childRes1: E, childRes2: E) -> E:
            return (childRes1[0] + childRes2[0], childRes1[1] + childRes2[1])

        def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
            """direction: 0: cur -> parent, 1: parent -> cur"""
            from_ = cur if direction == 0 else parent
            to_ = parent if direction == 0 else cur
            isPrime = int(P.isPrime(from_ + 1))  # !注意统计的是from
            zero, one = fromRes
            return (0, zero + 1) if isPrime else (zero + 1, one)

        R = Rerooting(n)
        for u, v in edges:
            R.addEdge(u - 1, v - 1)
        dp = R.rerooting(e=e, op=op, composition=composition, root=0)  # !答案不包含当前根节点
        res = 0
        for root in range(n):
            if P.isPrime(root + 1):
                res += dp[root][0]
            else:
                res += dp[root][1]
        return res // 2


# [[1,2],[1,3],[2,4],[2,5]]
if __name__ == "__main__":
    print(Solution().countPaths(5, [[1, 2], [1, 3], [2, 4], [2, 5]]))
