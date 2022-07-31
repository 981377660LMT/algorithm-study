# n<=100 稠密图
# 最小割拆点
from collections import defaultdict
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


def main() -> None:
    n, m = map(int, input().split())
    adjMap = defaultdict(set)
    costs = list(map(int, input().split()))


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
