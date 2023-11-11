# 完全二分木 顶点为1,2,...,inf
# n<=1e6
# x<=1e18
# !求n次移动后的坐标

# 直接移动会爆long long ，用python跑会TLE
# !这里需要用U来消除L/R
# !相邻元素消除用栈

import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, cur = map(int, input().split())
    s = input()  # U/R/L
    stack = []
    for char in s:
        if char == "U":
            if stack and stack[-1] != "U":
                stack.pop()
            else:
                stack.append(char)
        else:
            stack.append(char)

    for char in stack:
        if char == "U":
            cur //= 2
        elif char == "L":
            cur *= 2
        else:
            cur = cur * 2 + 1

    print(cur)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
