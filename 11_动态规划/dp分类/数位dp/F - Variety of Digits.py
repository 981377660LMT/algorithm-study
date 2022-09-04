# ABC235F
# !询问1到n中,数位包含c1，…， Cm的数的和。
# n<=1e(1e4)
# m<=10
# 0<=ci<...<cm<=9

from functools import lru_cache, reduce
from typing import Tuple
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


n = int(input())
m = int(input())
needs = list(map(int, input().split()))
target = reduce(lambda x, y: x | 1 << y, needs, 0)
pow10 = [1]
for i in range(1, len(str(n)) + 10):
    pow10.append(pow10[-1] * 10 % MOD)


def cal(upper: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, visited: int, hasLeadingZero: bool, isLimit: bool) -> Tuple[int, int]:
        """当前在第pos位,包含visited,hasLeadingZero表示有前导0,isLimit表示是否贴合上界"""
        if pos == 0:
            if visited & target == target:
                return 1, 0
            return 0, 0

        res = [0, 0]
        up = nums[pos - 1] if isLimit else 9
        for cur in range(up + 1):
            next = (0, 0)
            if hasLeadingZero and cur == 0:
                next = dfs(pos - 1, visited, True, (isLimit and cur == up))
            else:
                next = dfs(pos - 1, visited | (1 << cur), False, (isLimit and cur == up))
            res[0] = (res[0] + next[0]) % MOD
            res[1] = (res[1] + next[1] + next[0] * cur * pow10[pos - 1]) % MOD
        return tuple(res)

    nums = []
    while upper:
        div, mod = divmod(upper, 10)
        nums.append(mod)
        upper = div
    return dfs(len(nums), 0, True, True)[1]


print(cal(n))
