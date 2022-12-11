# C - Duodecim Ferra
# 十二铜表法

from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 長さ L の鉄の棒が東西方向に横たわっています。
# この棒を 11 箇所で切断して、12 本に分割します。このとき分割後の各棒の長さが全て正整数になるように分割しなければなりません。
# 分割のしかたが何通りあるかを求めてください。
# 二つの分割の方法は、一方で分割されているが他方で分割されていない位置が存在する場合に、そしてその場合に限って区別されます。

if __name__ == "__main__":

    L = int(input())

    @lru_cache(None)
    def dfs(index: int, remain: int, cur: int) -> int:
        if remain < 0:
            return 0
        if index == L:
            return 1 if (remain == 0) else 0
        return dfs(index + 1, remain, cur + 1) + dfs(index + 1, remain - 1, 1)

    res = dfs(1, 11, 1)
    dfs.cache_clear()
    print(res)
