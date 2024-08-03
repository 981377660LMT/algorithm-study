from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


mp = {"R": 0, "P": 1, "S": 2}
winBy = [1, 2, 0]

if __name__ == "__main__":
    N = int(input())
    S = input()
    nums = [mp[s] for s in S]
    winBys = [winBy[num] for num in nums]

    @lru_cache(None)
    def dfs(index: int, preType: int) -> int:
        if index == N:
            return 0
        res = 0
        for i in range(3):
            if i != preType and not winBy[i] == nums[index]:
                res = max(res, dfs(index + 1, i) + (winBys[index] == i))
        return res

    print(dfs(0, -1))
