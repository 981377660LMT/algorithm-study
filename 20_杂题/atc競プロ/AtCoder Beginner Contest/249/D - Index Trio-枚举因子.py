"""
1<=ai<=2e5
1<=n<=2e5
求(i,j,k)的对数 (1<=i,j,k<=n)
满足 ai/aj=ak

注意到ai只有2e5 因此可以存每个值有多少个
枚举因子 O(aloga)
"""

from collections import Counter
import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
MAX = int(2e5)


def main() -> None:
    n = int(input())
    nums = list(map(int, input().split()))
    counter = Counter(nums)
    res = 0
    for i in range(1, MAX + 1):
        for j in range(i, MAX + 1, i):
            res += counter[i] * counter[j] * counter[j // i]
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
