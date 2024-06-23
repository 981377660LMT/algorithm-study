from typing import List, Tuple
from collections import defaultdict
from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    N, K = map(int, input().split())
    S = input()

    nums = [0] * N
    for i in range(N):
        if S[i] == "A":
            nums[i] = 0
        elif S[i] == "B":
            nums[i] = 1
        else:
            nums[i] = 2

    def isPalindrom(preK):
        for i in range(len(preK) // 2):
            if preK[i] != preK[len(preK) - i - 1]:
                return False
        return True

    @lru_cache(None)
    def dfs(index: int, preK: Tuple) -> int:
        if index == N:
            return 1

        res = 0
        cur = nums[index]
        if cur == 0 or cur == 2:
            nextPreK = preK[1:] + (0,)
            if not isPalindrom(nextPreK):
                res += dfs(index + 1, nextPreK)
        if cur == 1 or cur == 2:
            nextPreK = preK[1:] + (1,)
            if not isPalindrom(nextPreK):
                res += dfs(index + 1, nextPreK)
        return res % MOD

    res = dfs(0, tuple(list(range(10, 10 + K))))
    print(res % MOD)
