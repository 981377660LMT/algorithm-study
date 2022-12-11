"""
求 1-n的排列个数, 且满足 m 个限制条件 [i,valueUpper,countUpper]
!A[1],A[2],...A[i]中不超过valueUpper的数的个数不超过countUpper

n<=18 m<=100
!全排列状压dp O(m*2^n)
"""

from functools import lru_cache
import sys
from typing import List, Tuple

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")


def countPermutation(n: int, restrictions: List[Tuple[int, int, int]]) -> int:
    @lru_cache(None)
    def dfs(index: int, visited: int) -> int:
        if index == n:
            return 1

        res = 0
        selected = [i + 1 for i in range(n) if visited & (1 << i)]
        for next in range(n):
            if visited & (1 << next):
                continue
            selected.append(next + 1)
            if check(index, selected):
                res += dfs(index + 1, visited | 1 << next)
            selected.pop()
        return res

    def check(index: int, nums: List[int]) -> bool:
        for valueUpper, countUpper in group[index]:
            count = 0
            for num in nums:
                if num <= valueUpper:
                    count += 1
                if count > countUpper:
                    return False
        return True

    group = [[] for _ in range(n)]  # !每个index处需要满足的限制条件
    for i, valueUpper, countUpper in restrictions:
        group[i].append((valueUpper, countUpper))
    res = dfs(0, 0)
    dfs.cache_clear()
    return res


if __name__ == "__main__":
    n, m = map(int, input().split())
    restrictions = []
    for _ in range(m):
        i, valueUpper, countUpper = map(int, input().split())
        i -= 1
        restrictions.append((i, valueUpper, countUpper))
    res = countPermutation(n, restrictions)
    print(res)
