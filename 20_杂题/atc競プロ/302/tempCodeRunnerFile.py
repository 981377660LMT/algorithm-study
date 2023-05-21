from math import ceil
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 体力が
# A の敵がいます。あなたは、
# 1 回攻撃をすることで敵の体力を
# B 減らすことが出来ます。

# 敵の体力を
# 0 以下にするためには、最小で何回攻撃をする必要があるでしょうか？
if __name__ == "__main__":
    A, B = map(int, input().split())
    print((A + B - 1) // B)
