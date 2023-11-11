from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    # 左右一列に N 個のマスが並んでおり、左から i 番目のマスの高さは H
    # i
    # ​
    #   です。

    # あなたは各マスについて 1 度ずつ次のいずれかの操作を行います。

    # マスの高さを 1 低くする。
    # 何もしない。
    # 操作をうまく行うことでマスの高さを左から右に向かって単調非減少にできるか求めてください。
    N = int(input())
    nums = list(map(int, input().split()))
