from collections import Counter, defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 0, 1, ? からなる文字列
# S および整数
# N が与えられます。
# S に含まれる ? をそれぞれ 0 または 1 に置き換えて
# 2 進整数とみなしたときに得られる値の集合を
# T とします。 たとえば、
# S= ?0? のとき、
# T={000
# (2)
# ​
#  ,001
# (2)
# ​
#  ,100
# (2)
# ​
#  ,101
# (2)
# ​
#  }={0,1,4,5} です。

# T に含まれる
# N 以下の値のうち最大のものを (
# 10 進整数として) 出力してください。
# N 以下の値が
# T に含まれない場合は、代わりに -1 を出力してください。
if __name__ == "__main__":
    s = input()
    n = int(input())

    # 数位dfs
    bin_ = list(map(int, bin(n)[2:]))
    bin_ = [0] * max(0, (len(s) - len(bin_))) + bin_

    def dfs(pos: int, isLimit: bool, cur: int):
        if pos == len(s):
            if not isLimit or cur <= n:
                yield cur
            return

        if s[pos] == "?":
            upper = bin_[pos] if isLimit else 1
            for x in range(upper, -1, -1):
                yield from dfs(pos + 1, isLimit and x >= upper, cur * 2 + x)
        else:
            x = int(s[pos])
            yield from dfs(pos + 1, isLimit and x >= bin_[pos], cur * 2 + x)

    res = next(dfs(0, True, 0), -1)
    print(res)
