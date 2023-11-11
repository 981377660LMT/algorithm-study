# 给出n个目标的位置(各不相同)
# 每次向x扔飞镖 1/3的概率飞到x-1,x,x+1
# !制定最优战略时(不会打没有目标的位置)，求击中所有目标所需回数的期望值
# n<=16
# 0<=xi<=15

# !状压+概率dp 答案只与当前剩下的目标状态有关
# !dp[state] = 1 + (dp[state1]/3 + dp[state2]/3 + dp[state3]/3)
# !注意如果nextState和state一样的话 需要合并同类项dp[state]

from functools import lru_cache, reduce
import sys


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
INF = int(4e18)

n = int(input())
places = sorted(map(int, input().split()))
target = reduce(lambda pre, cur: pre | (1 << cur), places, 0)


@lru_cache(None)
def dfs(remain: int) -> float:
    if remain == 0:
        return 0

    res = INF
    for i in range(16):
        remain1 = remain if i == 0 else remain & ~(1 << (i - 1))
        remain2 = remain & ~(1 << i)
        remain3 = remain if i == 15 else remain & ~(1 << (i + 1))  #  remain & ~(1 << (i + 1))

        step = 1
        p = 0
        if remain1 != remain:
            step += dfs(remain1) / 3
            p += 1 / 3
        if remain2 != remain:
            step += dfs(remain2) / 3
            p += 1 / 3
        if remain3 != remain:
            step += dfs(remain3) / 3
            p += 1 / 3

        if p != 0:
            res = min(res, step / p)

    return res


res = dfs(target)
dfs.cache_clear()
print(res)
