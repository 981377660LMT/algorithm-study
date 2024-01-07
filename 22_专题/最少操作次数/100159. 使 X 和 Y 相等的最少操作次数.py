# 100159. 使 X 和 Y 相等的最少操作次数
# https://leetcode.cn/problems/minimum-number-of-operations-to-make-x-and-target-equal/description/
# 给你两个正整数 x 和 target 。
# 一次操作中，你可以执行以下四种操作之一：
# 如果 x 是 11 的倍数，将 x 除以 11 。
# 如果 x 是 5 的倍数，将 x 除以 5 。
# 将 x 减 1 。
# 将 x 加 1 。
# 请你返回让 x 和 target 相等的 最少 操作次数。

# !除法对 x的影响远大于加减
# O(log(n)解法


from functools import lru_cache


class Solution:
    def minimumOperationsToMakeEqual(self, x: int, target: int) -> int:
        if x == target:
            return 0
        if x < target:
            return target - x

        @lru_cache(None)
        def dfs(cur: int) -> int:
            if cur <= target:
                return target - cur
            res1 = dfs(cur // 11) + 1 + cur % 11  # down
            res2 = dfs(cur // 11 + 1) + 1 + 11 - cur % 11  # up
            res3 = dfs(cur // 5) + 1 + cur % 5  # down
            res4 = dfs(cur // 5 + 1) + 1 + 5 - cur % 5  # up
            res5 = cur - target
            return min(res1, res2, res3, res4, res5)

        res = dfs(x)
        dfs.cache_clear()
        return res


if __name__ == "__main__":
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
