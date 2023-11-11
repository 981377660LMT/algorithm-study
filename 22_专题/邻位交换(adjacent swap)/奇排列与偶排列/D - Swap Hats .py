# 给定两个“RGB”这个字符串的排列A, B ，每次操作可以交换两个位置。
# 问能否正好操作1e18次把A 变成B。


# !奇排列与偶排列 一次交换会将奇排列变为偶排列
import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)

G1 = [("R", "G", "B"), ("G", "B", "R"), ("B", "R", "G")]


def main() -> None:
    s0, s1, s2 = input().split()
    t0, t1, t2 = input().split()
    flag = ((s0, s1, s2) in G1) ^ ((t0, t1, t2) in G1)
    print("No" if flag else "Yes")


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
