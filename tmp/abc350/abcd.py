from collections import defaultdict
from decimal import Decimal
from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 整数
# N が与えられます。あなたは次の
# 2 種類の操作を行うことができます。

# X 円払う。
# N を
# ⌊
# A
# N
# ​
#  ⌋ に置き換える。
# Y 円払う。
# 1 以上
# 6 以下の整数が等確率で出るサイコロを振る。その出目を
# b としたとき、
# N を
# ⌊
# b
# N
# ​
#  ⌋ に置き換える。
# ここで
# ⌊s⌋ は
# s 以下の最大の整数を表します。例えば
# ⌊3⌋=3、
# ⌊2.5⌋=2 です。

# 適切に操作を選択したとき、
# N を
# 0 にするまでに払う必要がある金額の期待値の最小値を求めてください。
# なお、サイコロの出目は操作ごとに独立であり、操作の選択はそれまでの操作の結果を確認してから行うことができます。


if __name__ == "__main__":
    N, A, X, Y = map(int, input().split())
    F = Decimal(25)

    @lru_cache(None)
    def dfs(cur: int) -> float:
        if cur == 0:
            return 0

        # X 円払う
        res = Decimal(X + dfs(cur // A))
        # Y 円払う
        res2 = 0
        for i in range(2, 7):
            res2 += Decimal((6 * Y + 5 * dfs(cur // i)))
        res2 /= F
        return min(res, res2)

    print(dfs(N))
