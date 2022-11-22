# 在 "100 game" 这个游戏中，两名玩家轮流选择从 1 到 10 的任意整数，累计整数和，
# 先使得累计整数和 达到或超过  100 的玩家，即为胜者。

# 如果我们将游戏规则改为 “玩家 不能 重复使用整数” 呢？

# 例如，两个玩家可以轮流从公共整数池中抽取从 1 到 15 的整数（不放回），直到累计整数和 >= 100。

# 给定两个整数 maxChoosableInteger （整数池中可选择的最大数）和 desiredTotal（累计和），
# !若先出手的玩家是否能稳赢则返回 true ，否则返回 false 。假设两位玩家游戏时都表现 最佳 。


from functools import lru_cache


class Solution:
    def canIWin(self, maxNum: int, target: int) -> bool:
        @lru_cache(None)
        def dfs(curSum: int, visited: int) -> bool:
            for i in range(1, maxNum + 1):
                if 1 << i & visited:
                    continue
                # 自己赢 = 自己赢或对手不赢
                if curSum + i >= target or not dfs(curSum + i, visited | 1 << i):
                    return True
            return False

        if (maxNum + 1) * maxNum // 2 < target:
            return False

        res = dfs(0, 0)
        dfs.cache_clear()
        return res
