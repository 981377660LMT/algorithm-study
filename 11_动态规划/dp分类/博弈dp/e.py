from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 高橋君と青木君がすごろくをします。
# 高橋君ははじめ地点
# A、青木君ははじめ地点
# B にいて、交互にサイコロを振ります。
# 高橋君が振るサイコロは
# 1,2,…,P の出目が一様ランダムに出るサイコロで、青木君が振るサイコロは
# 1,2,…,Q の出目が一様ランダムに出るサイコロです。
# 地点
# x にいるときに自分の振ったサイコロの出目が
# i であるとき、地点
# min(x+i,N) に進みます。
# 地点
# N に先に着いた人をすごろくの勝者とします。
# 高橋君が先にサイコロを振るとき、高橋君が勝つ確率を
# mod 998244353 で求めてください。
if __name__ == "__main__":
    N, A, B, P, Q = map(int, input().split())
    invP = pow(P, MOD - 2, MOD)
    invQ = pow(Q, MOD - 2, MOD)

    @lru_cache(None)
    def dfs(pos1: int, pos2: int, turn: int) -> int:
        if turn == 0:  # takahashi
            if pos1 >= N:
                return 1
            res = 0
            for i in range(1, P + 1):
                nextPos = min(pos1 + i, N)
                if nextPos == N:
                    res += invP
                else:
                    res += (1 - dfs(nextPos, pos2, 1)) * invP % MOD
            return res % MOD
        else:
            res = 0
            for i in range(1, Q + 1):
                nextPos = min(pos2 + i, N)
                if nextPos == N:
                    res += invQ
                else:
                    res += (1 - dfs(pos1, nextPos, 0)) * invQ % MOD
            return res % MOD

    res = dfs(A, B, 0)
    print(res)
