# https://atcoder.jp/contests/abc235/tasks/abc235_g
# 三种颜色的球分别有A,B,C个，同种颜色的球不区分，现将其放入N个有标号的箱子内，要求
# - 每个箱子至少有1个球。
# - 每个箱子内至多有1个同种球
# - 不必放完所有的球。

# 求方案数模998244353,1<N≤5e6, 0≤A, B,C ≤N。
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# TODO
