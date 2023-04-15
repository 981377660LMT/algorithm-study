from heapq import heapify, heappop, heappush
from typing import List


class Solution:
    def nthSuperUglyNumber(self, k: int, primes: List[int]) -> int:
        """超级丑数 是一个正整数，并满足其所有质因数都出现在质数数组 primes 中。

        给你一个整数 k 和一个整数数组 primes ，返回第 k 个 超级丑数 。
        1<=k<=1e5 len(primes)<=100
        """
        n = len(primes)
        pq = [(primes[i], i, 0) for i in range(n)]
        heapify(pq)

        res = [1]
        while len(res) < k:
            val, row, col = heappop(pq)
            if val != res[-1]:
                res.append(val)
            nextVal = res[col + 1] * primes[row]
            heappush(pq, (nextVal, row, col + 1))

        return res[-1]


print(Solution().nthSuperUglyNumber(k=12, primes=[2, 7, 13, 19]))
