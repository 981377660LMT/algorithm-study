from functools import lru_cache


# A - Pay to Win
# https://atcoder.jp/contests/agc044/tasks/agc044_a
# 给定一个初始为0的数，每次操作可以：
# 1. 花费A,将数变为2x
# 2. 花费B,将数变为3x
# 3. 花费C,将数变为5x
# 4. 花费D,将数变为x+1 或者 x-1
# 问最少花费多少可以将数变为N(n<=1e18)
#
# !等价从 N 开始变到 0
def payToWin() -> None:
    @lru_cache(None)
    def dfs(cur: int) -> int:
        if cur == 0:
            return 0
        if cur == 1:
            return D

        res = D * cur

        div, mod = cur // 2, cur % 2
        if mod == 0:
            res = min(res, A + dfs(div))
        else:
            res = min(res, A + D + dfs(div), A + D + dfs(div + 1))  # ! +1 or -1

        div, mod = cur // 3, cur % 3
        if mod == 0:
            res = min(res, B + dfs(div))
        elif mod == 1:
            res = min(res, B + D + dfs(div))  # ! -1
        else:
            res = min(res, B + D + dfs(div + 1))  # ! +1

        div, mod = cur // 5, cur % 5
        if mod == 0:
            res = min(res, C + dfs(div))
        elif mod == 1:
            res = min(res, C + D + dfs(div))  # ! -1
        elif mod == 2:
            res = min(res, C + D + D + dfs(div))  # ! -2
        elif mod == 3:
            res = min(res, C + D + D + dfs(div + 1))  # ! +2
        else:
            res = min(res, C + D + dfs(div + 1))  # ! +1

        return res

    T = int(input())
    for _ in range(T):
        N, A, B, C, D = map(int, input().split())

        dfs.cache_clear()
        print(dfs(N))


payToWin()
