from functools import lru_cache

# 给定两个整数 maxChoosableInteger （整数池中可选择的最大数）和 desiredTotal（累计和），
# 若先出手的玩家是否能稳赢则返回 true ，否则返回 false 。假设两位玩家游戏时都表现 最佳 。


class Solution:
    def canIWin(self, upper: int, target: int) -> bool:
        """2^n*n 每次dfs不用重新计算cur，因为curSum由visted唯一确定，两种的时间复杂度都是 O(2^n*n)"""

        @lru_cache(None)
        def dfs(curSum: int, visited: int) -> bool:
            if curSum >= target:
                return True

            for select in range(1, upper + 1):
                if (visited >> select) & 1:
                    continue
                # 自己赢=自己赢或对手不赢
                if curSum + select >= target or not dfs(curSum + select, visited | (1 << select)):
                    return True
            return False

        if upper * (upper + 1) // 2 < target:
            return False

        res = dfs(0, 0)
        dfs.cache_clear()
        return res

    def canIWin2(self, upper: int, target: int) -> bool:
        """2^n*n 会慢一些"""

        @lru_cache(None)
        def dfs(state: int) -> bool:
            curSum = sum(i for i in range(1, upper + 1) if (state >> i) & 1)
            for select in range(1, upper + 1):
                if (state >> select) & 1:
                    continue
                # 自己赢=自己赢或对手不赢
                if curSum + select >= target or not dfs(curSum + select, state | (1 << select)):
                    return True
            return False

        if upper * (upper + 1) // 2 < target:
            return False

        res = dfs(0)
        dfs.cache_clear()
        return res

