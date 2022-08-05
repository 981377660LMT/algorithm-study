# 定义f(a)是和a的位数相同且小于等于的正整数的个数f(1)= 1,f(2)=2,f(10)= 1 ...
# 给定一个a，求f(1) + f(2) +f(3)＋ ...f(n)的值模998244353的值
# n<=1e18

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def cal(length: int, upper: int) -> int:
    """长度为length,上限为upper时这一段的和"""
    start, end = 10 ** (length - 1), min(10**length - 1, upper)
    first, last = 1, end - start + 1
    count = end - start + 1
    return (first + last) * count // 2


def main() -> None:
    n = int(input())
    k = len(str(n))
    res = 0
    for i in range(1, k + 1):
        res += cal(i, n)
        res %= MOD
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
