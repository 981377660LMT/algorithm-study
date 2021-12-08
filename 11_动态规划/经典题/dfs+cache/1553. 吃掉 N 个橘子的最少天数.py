from functools import lru_cache

# 吃掉一个橘子。
# 如果剩余橘子数 n 能被 2 整除，那么你可以吃掉 n/2 个橘子。
# 如果剩余橘子数 n 能被 3 整除，那么你可以吃掉 2*(n/3) 个橘子。


class Solution:
    def minDays(self, n: int) -> int:
        @lru_cache(None)
        def dfs(n: int) -> int:
            if n <= 1:
                return n
            return 1 + min(dfs(n // 2) + n % 2, dfs(n // 3) + n % 3)

        return dfs(n)

    # slow
    def minDays2(self, n: int) -> int:
        res = 0
        queue = [n]
        visited = set()
        while queue:  # bfs
            nextQueue = []
            for x in queue:
                if x == 0:
                    return res
                visited.add(x)
                if x - 1 not in visited:
                    nextQueue.append(x - 1)
                if x % 2 == 0 and x // 2 not in visited:
                    nextQueue.append(x // 2)
                if x % 3 == 0 and x // 3 not in visited:
                    nextQueue.append(x // 3)
            res += 1
            queue = nextQueue


print(Solution().minDays(10))
# 输出：4
# 解释：你总共有 10 个橘子。
# 第 1 天：吃 1 个橘子，剩余橘子数 10 - 1 = 9。
# 第 2 天：吃 6 个橘子，剩余橘子数 9 - 2*(9/3) = 9 - 6 = 3。（9 可以被 3 整除）
# 第 3 天：吃 2 个橘子，剩余橘子数 3 - 2*(3/3) = 3 - 2 = 1。
# 第 4 天：吃掉最后 1 个橘子，剩余橘子数 1 - 1 = 0。
# 你需要至少 4 天吃掉 10 个橘子。
