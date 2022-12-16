from collections import Counter
from itertools import accumulate
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
if __name__ == "__main__":
    import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
if __name__ == "__main__":
    #   高橋君は AtCoder のコンテストに参加しています。

    # このコンテストでは、 N 問の問題が出題されます。

    # 高橋君はコンテスト中に M 回の提出を行いました。

    # i 回目の提出は p
    # i
    # ​
    #   番目の問題への提出であり、結果は S
    # i
    # ​
    #   (AC または WA) でした。

    # 高橋君の正答数は、AC を 1 回以上出した問題の数です。

    # 高橋君のペナルティ数は、高橋君が AC を 1 回以上出した各問題において、初めて AC を出すまでに出した WA の数の総和です。

    # 高橋君の正答数とペナルティ数を答えてください。
    n, m = map(int, input().split())
    pairs = [list(input().split()) for _ in range(m)]
