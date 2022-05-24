# https://www.acwing.com/solution/content/75714/
# 它休息的这 B 个小时不一定连续，可以分成若干段，但是在每段的第一个小时，它需要从清醒逐渐入睡，不能恢复体力，从下一个小时开始才能睡着。
from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
# n<=3000

# 分情况：
@lru_cache(None)
def dfs1(index: int, count: int, hasPre: bool) -> int:
    """第n天睡觉"""
    if count > need:
        return -int(1e10)
    if index == n - 1:
        if count != need:
            return -int(1e10)
        return scores[index] if hasPre else 0
    res = -int(1e20)
    # 不睡
    res = max(res, dfs1(index + 1, count, False))
    # 睡
    res = max(res, dfs1(index + 1, count + 1, True) + (scores[index] if hasPre else 0))
    return res


@lru_cache(None)
def dfs2(index: int, count: int, hasPre: bool) -> int:
    """第n天不睡觉"""
    if count > need:
        return -int(1e10)
    if index == n - 1:
        if count != need:
            return -int(1e10)
        return 0

    res = -int(1e20)
    # 不睡
    res = max(res, dfs2(index + 1, count, False))
    # 睡
    res = max(res, dfs2(index + 1, count + 1, True) + (scores[index] if hasPre else 0))
    return res


n, need = map(int, input().split())
scores = []
for _ in range(n):
    scores.append(int(input()))

print(max(dfs1(0, 1, True), dfs2(0, 0, False)))


# [2,0,3,1,4]
# 这头牛每天 3 点入睡，睡到次日 1 点，即 [1,4,2] 时间段休息，每天恢复体力值最大，为 0+4+2=6。
