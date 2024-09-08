# No.794 チーム戦 (2)
# 给定长为偶数n的数组nums，求和不超过k的二元对数的分配方案数。
# https://yukicoder.me/problems/no/794
#

from bisect import bisect_right
import sys


input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def min2(a: int, b: int) -> int:
    return a if a < b else b


if __name__ == "__main__":
    n, k = map(int, input().split())
    nums = list(map(int, input().split()))

    nums.sort()
    res = 1
    for i in range(n // 2):
        pos = bisect_right(nums, k - nums[-i - 1])
        res *= min2(n - i - 1, pos) - i
        res %= MOD

    print(res)
