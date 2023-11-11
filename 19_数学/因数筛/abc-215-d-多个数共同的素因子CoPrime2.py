"""
给n个数,求在1-m中所有与这n个数互质的所有数。
n,m,ai<=1e5
"""


from typing import List
from prime import EratosthenesSieve

E = EratosthenesSieve(int(1e5 + 10))


def getCoPrimes(nums: List[int], upper: int) -> List[int]:
    """求1-upper中与nums中所有数互质的数 O(nlogn)"""
    primes = set()
    for num in nums:
        primes |= set(E.getPrimeFactors(num))
    visited = [False] * (upper + 1)
    for p in primes:
        for multi in range(p, upper + 1, p):
            visited[multi] = True
    return [num for num in range(1, upper + 1) if not visited[num]]


############################################################
import sys


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    n, m = map(int, input().split())
    nums = list(map(int, input().split()))
    res = getCoPrimes(nums, m)
    print(len(res))
    print(*res, sep="\n")
