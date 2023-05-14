from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 英大文字・英小文字からなる長さ
# 4 の文字列で、以下の
# 2 条件をともに満たすものを DDoS 型文字列と呼ぶことにします。

# 1,2,4 文字目が英大文字で、
# 3 文字目が英小文字である
# 1,2 文字目が等しい
# 例えば DDoS, AAaA は DDoS 型文字列であり、ddos, IPoE は DDoS 型文字列ではありません。

# 英大文字・英小文字および ? からなる文字列
# S が与えられます。
# S に含まれる ? を独立に英大文字・英小文字に置き換えてできる文字列は、
# S に含まれる ? の個数を
# q として
# 52
# q
#   通りあります。 このうち DDoS 型文字列を部分列に含まないものの個数を
# 998244353 で割ったあまりを求めてください。
if __name__ == "__main__":
    s = input()

    @lru_cache(None)
    def dfs(index: int, D2: bool, small: bool, big: bool) -> int:
        if index == len(s):
            return 1
