import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# N 個の実数の区間が与えられます。
# i(1≤i≤N) 番目の区間は
# [l
# i
# ​
#  ,r
# i
# ​
#  ] です。
# i 番目の区間と
# j 番目の区間が共通部分を持つような組
# (i,j)(1≤i<j≤N) の個数を求めてください。
if __name__ == "__main__":
    N = int(input())
    intervals = [list(map(int, input().split())) for _ in range(N)]
    intervals.sort(key=lambda x: x[0])
    ans = 0
