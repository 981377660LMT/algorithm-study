# 幅 A、奥行き B、高さ C の直方体の形をしたケーキがあります。
# あなたはケーキに対して、次の操作を行うことができます。
# - ある面に平行な方向に切断する
# - ただし、ケーキは動かしてはならない（複数のケーキに分割されている場合、これらを変形したり、別々に切ることはできない）
# !最小何回の操作で、全てのピースを立方体にすることができますか？
# 1≤A,B,C≤1e18

from math import gcd
import sys

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)

A, B, C = map(int, input().split())
gcd_ = gcd(A, gcd(B, C))  # version < 3.9

# 每个方向切几刀
res1, res2, res3 = A // gcd_ - 1, B // gcd_ - 1, C // gcd_ - 1
print(res1 + res2 + res3)

