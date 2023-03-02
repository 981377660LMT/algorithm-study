from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 1 から
# N までの番号がついた
# N 枚のカードが一列に並んでいて、各
# i (1≤i<N) に対してカード
# i とカード
# i+1 が隣り合っています。 カード
# i の表には
# A
# i
# ​
#   が、裏には
# B
# i
# ​
#   が書かれており、最初全てのカードは表を向いています。

# 今から、
# N 枚のカードのうち好きな枚数 (
# 0 枚でも良い) を選んで裏返すことを考えます。 裏返すカードの選び方は
# 2
# N
#   通りありますが、そのうち以下の条件を満たすものの数を
# 998244353 で割った余りを求めてください。

# 選んだカードを裏返した後、どの隣り合う
# 2 枚のカードについても、向いている面に書かれた数が相異なる。

if __name__ == "__main__":
    n = int(input())
    A, B = [], []
    for _ in range(n):
        a, b = input().split()
        A.append(a)
        B.append(b)

    @lru_cache(None)
    def dfs(index: int, preFlip: bool) -> int:
        if index == n:
            return 1
        if index == 0:
            return dfs(index + 1, True) + dfs(index + 1, False)
        res = 0
        pre = B[index - 1] if preFlip else A[index - 1]
        if pre != A[index]:
            res += dfs(index + 1, False)
        if pre != B[index]:
            res += dfs(index + 1, True)
        return res % MOD

    res = dfs(0, False)
    print(res)
