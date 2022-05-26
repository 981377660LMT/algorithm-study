class Solution:
    def maxA(self, n: int) -> int:
        """肯定是从某个时间开始不断按v"""

        @lru_cache(None)
        def dfs(cur: int, step: int, copy: int) -> int:
            if step >= n:
                return cur if step == 0 else -int(1e20)

            res1 = dfs(cur + max(1, copy), step + 1, copy)  # 打字或粘贴
            res2 = dfs(cur, step + 2, cur)
            return max(res1, res2)

        return dfs(0, 0, 0)