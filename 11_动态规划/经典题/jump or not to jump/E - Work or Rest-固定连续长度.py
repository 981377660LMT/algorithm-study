# 在某个星球上，一周有 n 天，
# 你作为这个星球的国王，需要给一周的每一天都定为休息日或者工作日，
# 至少要有一天被定为休息日。

# 当某一天为休息日时，人们的生产量为 0
# 当某一天为工作日时，人们的生产量为 scores[min(a,b)]
# a 为这一天距离上一个休息日的天数，
# b 为这一天距离下一个休息日的天数；
# !请问人们每周生产量的最大值为多少。
# !n<=5000


# !不失一般性,设每周的第一天为休息日
# 注意到连续工作1天, 得分为scores[0]
# 连续工作2天, 得分为scores[0] + scores[0]
# 连续工作3天, 得分为scores[0] + scores[1] + scores[0]
# !dp[i]表示前i天的最大生产量,转移时枚举连续工作j天
# !(固定连续长度从而解决了环的问题)

from functools import lru_cache
from typing import List


def workOrRest(n: int, scores: List[int]) -> int:
    @lru_cache(None)
    def dfs(index: int) -> int:
        if index == n:
            return 0
        res = 0
        for work in range(n + 1):  # 工作work天后休息1天
            if index + work + 1 > n:
                break
            res = max(res, preSum[work] + dfs(index + work + 1))
        return res

    preSum = [0] * (n + 1)  # 连续工作j天的得分
    for i in range(1, n + 1):
        preSum[i] = preSum[i - 1] + scores[(i - 1) // 2]

    res = dfs(0)
    dfs.cache_clear()
    return res


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    scores = list(map(int, input().split()))
    print(workOrRest(n, scores))
