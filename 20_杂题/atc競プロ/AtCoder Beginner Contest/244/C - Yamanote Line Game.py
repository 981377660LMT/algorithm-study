# 报数字游戏
# n<=1000
# 可用的数字为 1到 2*n+1
# 两人互相报数字 不能重复报 报不了的就输
# !这个游戏先手必胜 请根据队友的报数来返回你的报数

# !注意到数字可以两两配对 和为 2*n+2 (除去n+1)

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n = int(input())
    pairSum = 2 * n + 2
    print(n + 1, flush=True)
    for _ in range(n):
        num = int(input())
        print(pairSum - num, flush=True)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
