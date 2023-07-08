# 给你两个长度为 n 下标从 0 开始的整数数组 cost 和 time ，
# 分别表示给 n 堵不同的墙刷油漆需要的开销和时间。你有两名油漆匠：

# 一位需要 付费 的油漆匠，刷第 i 堵墙需要花费 time[i] 单位的时间，开销为 cost[i] 单位的钱。
# 一位 免费 的油漆匠，刷 任意 一堵墙的时间为 1 单位，开销为 0 。
# 但是必须在付费油漆匠 工作 时，免费油漆匠才会工作。
# 请你返回刷完 n 堵墙最少开销为多少。


# 1 <= cost.length <= 500
# cost.length == time.length
# 1 <= cost[i] <= 1e6
# 1 <= time[i] <= 500

# 没做出来的原因是看到time[i]为500，把time当状态了，做不出来(其实也可以，就是需要对状态剪枝)
# !dp[i][diff]表示前i个墙壁，付费油漆匠时间-免费油漆匠刷的个数为diff时的最小开销


from functools import lru_cache
from typing import List

INF = int(1e18)


def min(x, y):
    if x < y:
        return x
    return y


class Solution:
    def paintWalls(self, cost: List[int], time: List[int]) -> int:
        @lru_cache(None)
        def dfs(index: int, diff: int) -> int:
            if diff >= (n - index):
                return 0
            if index == n:
                return 0 if diff >= 0 else INF
            res1 = dfs(index + 1, diff + time[index])
            res2 = dfs(index + 1, diff - 1)
            return min(res1 + cost[index], res2)

        n = len(cost)
        res = dfs(0, 0)
        dfs.cache_clear()
        return res
