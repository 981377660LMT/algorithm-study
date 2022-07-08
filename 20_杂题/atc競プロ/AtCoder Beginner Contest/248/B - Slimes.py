from math import ceil, log
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)

EPS = 1e-6


def main() -> None:
    a, b, k = map(int, input().split())
    # すぬけくんが 1 回叫ぶたびに、スライムは K 倍に増殖します。
    # スライムが B 匹以上(>=B)になるには、すぬけくんは最小で何回叫ぶ必要があるでしょうか
    print(ceil(log((b / a), k)))  # !注意log函数先传指数，再传底数


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
