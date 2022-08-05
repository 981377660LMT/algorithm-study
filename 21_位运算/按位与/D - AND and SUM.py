# 求有多少组非負整数(x,y) 使得
# x&y=a x+y=s
# T<=1e5
# a,s<=2^60

# !注意公式 (a & b) << 1 + (a ^ b) = a + b
# 两数之和=两数与*2+两数异或
# !则可以求出两个数的与和两个数的异或 进而按位求解
# !两个数的与和两个数的异或在每一位上可能是 (0,0),(1,0),(0,1)

from collections import defaultdict
from itertools import product
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


MAPPING = defaultdict(int)
for a, b in product((0, 1), repeat=2):
    MAPPING[(a & b, a ^ b)] += 1


def work():
    and_, sum_ = map(int, input().split())
    xor_ = sum_ - and_ * 2
    res = 1
    if xor_ < 0:  # 负数情况特判
        print("No")
        return

    for i in range(63):
        bit1, bit2 = (and_ >> i) & 1, (xor_ >> i) & 1
        res *= MAPPING[(bit1, bit2)]  # type: ignore
    print(["No", "Yes"][res > 0])


def main() -> None:
    T = int(input())
    for _ in range(T):
        work()


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
