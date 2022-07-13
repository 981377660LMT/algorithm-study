# C - XX to XXX
# 给定两个串S和T，每次可以向s中相邻且相同的两个字符中间塞一个相同的字符。
# 问若干次操作后S是否能变成T。

# !groupby 分组
from itertools import groupby
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    s = input()
    t = input()
    group1 = [(char, len(list(group))) for char, group in groupby(s)]
    group2 = [(char, len(list(group))) for char, group in groupby(t)]
    if len(group1) != len(group2):
        print("No")
        exit(0)
    for (c1, n1), (c2, n2) in zip(group1, group2):
        if c1 != c2 or n1 > n2 or (n1 == 1 and n2 != 1):
            print("No")
            exit(0)
    print("Yes")


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
