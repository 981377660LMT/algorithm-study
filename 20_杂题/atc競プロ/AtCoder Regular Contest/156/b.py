import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 有限個の非負整数からなる多重集合
# S にたいして、
# mex(S) を、
# S に含まれない最小の非負整数と定義します。例えば、
# mex({0,0,1,3})=2,mex({1})=0,mex({})=0 です。

# 黒板に
# N 個の非負整数が書かれており、
# i 番目の非負整数は
# A
# i
# ​
#   です。

# あなたは、以下の操作をちょうど
# K 回行います。

# 黒板に書かれている非負整数を
# 0 個以上選ぶ。選んだ非負整数からなる多重集合を
# S として、
# mex(S) を黒板に
# 1 個書き込む。
# 最終的に黒板に書かれている非負整数の多重集合としてありうるものの個数を
# 998244353 で割ったあまりを求めてください。
if __name__ == "__main__":
    n, k = map(int, input().split())
    nums = list(map(int, input().split()))

    # 枚举mex
