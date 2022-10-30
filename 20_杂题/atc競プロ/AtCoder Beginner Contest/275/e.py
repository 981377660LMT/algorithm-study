from functools import lru_cache
import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 高橋君は双六で遊んでいます。
# この双六には 0 から N の番号がついた N+1 個のマスがあります。 高橋君はマス 0 からスタートし、マス N を目指します。
# この双六では、1 から M までの M 種類の目が等確率で出るルーレットを使います。 高橋君はルーレットを回して出た目の数だけ進みます。もし、マス N を超えて進むことになる場合、マス N を超えた分だけ引き返します。
# 例えば、 N=4 で高橋君がマス 3 にいるとき、ルーレットを回して出た目が 4 の場合は、マス 4 を 3+4−4=3 マス超えてしまいます。そのため、 3 マスだけマス 4 から引き返し、マス 1 に移動します。
# 高橋君がマス N に到達するとゴールとなり、双六を終了します。
# 高橋君がルーレットを K 回まで回す時、ゴールできる確率を mod 998244353 で求めてください。


if __name__ == "__main__":

    @lru_cache(None)
    def dfs(pos: int, remain: int) -> int:
        if pos == n:
            return pow(m, remain, MOD)  # !剩下投掷次数都可以
        if remain == 0:
            return 0

        res = 0
        for cur in range(1, m + 1):
            next = pos + cur
            if next > n:
                overflow = next - n
                next = n - overflow
            res += dfs(next, remain - 1)
            res %= MOD
        return res

    n, m, k = map(int, input().split())
    res = dfs(0, k)
    all_ = pow(m, k, MOD)
    print(res * pow(all_, MOD - 2, MOD) % MOD)
