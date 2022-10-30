from functools import lru_cache
import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 非負整数 x に対し定義される関数 f(x) は以下の条件を満たします。

# f(0)=1
# 任意の正整数 k に対し f(k)=f(⌊
# 2
# k
# ​
#  ⌋)+f(⌊
# 3
# k
# ​
#  ⌋)
# ここで、⌊A⌋ は A の小数点以下を切り捨てた値を指します。

# このとき、 f(N) を求めてください。
if __name__ == "__main__":
    n = int(input())

    @lru_cache(None)
    def dfs(cur: int) -> int:
        if cur == 0:
            return 1
        return dfs(cur // 2) + dfs(cur // 3)

    print(dfs(n))
