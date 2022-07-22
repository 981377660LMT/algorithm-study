# 有v升洗发水，a,b,c轮着用，每个人分别用a,b,c升。轮到谁洗发水不够
# a不够输出F,b输出M,c输出T

# !周期问题 注意一开始模一下

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    remain, a, b, c = map(int, input().split())
    sum_ = a + b + c
    remain %= sum_
    for cost, cand in zip((a, b, c), ("F", "M", "T")):
        next = remain - cost
        if next < 0:
            print(cand)
            exit(0)
        remain = next


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
