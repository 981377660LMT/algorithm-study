from collections import Counter
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# N 枚の靴下があります。
# i 枚目の靴下の色は
# A
# i
# ​
#   です。

# あなたは以下の操作をできるだけ多い回数行いたいです。最大で何回行うことができますか？

# まだペアになっていない靴下の中から同じ色の靴下を
# 2 枚選んでペアにする。
if __name__ == "__main__":
    n = int(input())
    A = list(map(int, input().split()))
    counter = Counter(A)
    print(sum([v // 2 for v in counter.values()]))
