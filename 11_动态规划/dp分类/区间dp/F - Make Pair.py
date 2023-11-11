"""
有2N个学生.编号1,2,3...2N,其中有M对学生关系友好
教师每次可以从中选出一对相邻的学生,且他们关系友好
然后将这对学生删除(注意删除后学生的相邻关系改变).
求出删除所有学生的方案数,取模.
n<=200


https://www.cnblogs.com/dream1024/p/15244203.html
!区间dp
"""

from functools import lru_cache
import sys


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
a
fac = [1]
ifac = [1]
for i in range(1, int(1e3) + 10):
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


if __name__ == "__main__":
    n, m = map(int, input().split())
    pairs = []
    for _ in range(m):
        x, y = map(int, input().split())
        x, y = x - 1, y - 1
        pairs.append((x, y))

    # @lru_cache(None)
    # def dfs(left: int, right: int) -> int:
    #     """删除[left,right]的方案数"""
    #     if left >= right:
    #         return 0
    #     if left + 1 == right:
    #         return 1

    # print(dfs(0, 2 * n - 1))


# TODO
