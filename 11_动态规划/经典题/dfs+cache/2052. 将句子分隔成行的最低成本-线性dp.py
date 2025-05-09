from functools import lru_cache


INF = int(1e18)


class Solution:
    def minimumCost(self, sentence: str, k: int) -> int:
        words = sentence.split(" ")
        n = len(words)
        presum = [0] * (n + 1)
        for i, w in enumerate(words):
            presum[i + 1] = presum[i] + len(w)

        @lru_cache(None)
        def dfs(i: int):
            if presum[n] - presum[i] + n - 1 - i <= k:  # 最后一行
                return 0
            res, j = INF, i + 1
            while j < n and (nxt := presum[j] - presum[i] + j - 1 - i) <= k:
                res = min(res, (k - nxt) ** 2 + dfs(j))
                j += 1
            return res

        return dfs(0)


print(Solution().minimumCost(sentence="i love leetcode", k=12))
