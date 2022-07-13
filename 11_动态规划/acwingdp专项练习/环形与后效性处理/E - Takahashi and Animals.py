# 每次可以花费A[i] 来喂动物i和i+1 (取模)
# 所有动物都喂到的最小花费
# 环形分类:最后一个A[-1]选不选
# !2<=n<=3e5

from functools import lru_cache
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    @lru_cache(None)
    def dfs1(index: int, hasPre: int) -> int:
        """最后一个选了"""
        if index == n - 1:
            return nums[index]
        res = dfs1(index + 1, True) + nums[index]
        if hasPre:
            res = min(res, dfs1(index + 1, False))
        return res

    @lru_cache(None)
    def dfs2(index: int, hasPre: int) -> int:
        """最后一个没选"""
        if index == n - 1:
            return 0 if hasPre else int(1e18)
        res = dfs2(index + 1, True) + nums[index]
        if hasPre:
            res = min(res, dfs2(index + 1, False))
        return res

    n = int(input())
    nums = list(map(int, input().split()))
    print(min(dfs1(0, True), dfs2(0, False)))


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
