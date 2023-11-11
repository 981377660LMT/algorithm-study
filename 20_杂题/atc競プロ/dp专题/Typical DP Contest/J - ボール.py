# 给出n个目标的位置(各不相同)
# 每次向x扔飞镖 1/3的概率飞到x-1,x,x+1
# 制定最优战略时，求击中所有目标所需回数的期望值
# n<=16
# 0<=xi<=15

# !概率dp 答案只与当前剩下的目标有关
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
    for cur in range(16):  # 最优战略:不会扔的3个位置都没有目标
        if (
            (remain >> cur) & 1 == 0
            and (cur == 15 or (remain >> (cur + 1)) & 1 == 0)
            and (cur == 0 or (remain >> (cur - 1)) & 1 == 0)
        ):
            continue

        cand = 1
        count = 0
        for next in (cur - 1, cur, cur + 1):
            if next < 0 or next >= 16:
                continue
            if (remain >> next) & 1:
                count += 1
                cand += dfs(remain & ~(1 << next)) / 3
        res = min(res, cand * 3 / count)  # 有效的概率为count/3 因此要乘以3/count
    return res


res = dfs(target)
dfs.cache_clear()
print(res)
