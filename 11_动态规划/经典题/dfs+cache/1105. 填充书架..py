from typing import List
from functools import lru_cache

# 摆放书的顺序与你整理好的顺序相同。
# 以这种方式布置书架，返回书架整体可能的最小高度。
# 1 <= books.length <= 1000
# 1 <= thicknessi <= shelfWidth <= 1000
# 1 <= heighti <= 1000


INF = int(1e18)

# dp[i] = min(dp[j-1]+max(h[j],h[j+1],...,h[i])) 且 sum(w[j],w[j+1],...,w[i]) <= W


class Solution:
    def minHeightShelves(self, books: List[List[int]], shelfWidth: int) -> int:
        """O(n^2) dfs[i]表示前i本书的高度之和最小值"""

        @lru_cache(None)
        def dfs(index: int) -> int:
            # dp(index)返回装到第index本书的最小高度 内层循环看放哪些书到同一层
            if index == n:
                return 0

            curMax, remain = 0, shelfWidth
            res = INF
            # 连续放几本书
            for select in range(1, n - index + 1):
                width, height = books[index + select - 1]
                if width > remain:
                    break
                curMax = max(curMax, height)
                remain -= width
                res = min(res, curMax + dfs(index + select))
            return res

        n = len(books)
        res = dfs(0)
        dfs.cache_clear()
        return res

    def minHeightShelves2(self, books: List[List[int]], shelfWidth: int) -> int:
        """O(n^2) dp[i]表示前i本书的最小高度之和"""
        n = len(books)
        dp = [INF] * (n + 1)
        dp[0] = 0
        for i in range(n):
            curMax, remain = 0, shelfWidth
            for select in range(1, n - i + 1):
                width, height = books[i + select - 1]
                if width > remain:
                    break
                curMax = max(curMax, height)
                remain -= width
                dp[i + select] = min(dp[i + select], dp[i] + curMax)
        return dp[-1]


print(
    Solution().minHeightShelves(
        books=[[1, 1], [2, 3], [2, 3], [1, 1], [1, 1], [1, 1], [1, 2]], shelfWidth=4
    )
)
print(
    Solution().minHeightShelves2(
        books=[[1, 1], [2, 3], [2, 3], [1, 1], [1, 1], [1, 1], [1, 2]], shelfWidth=4
    )
)

# 输出：6
# 解释：
# 3 层书架的高度和为 1 + 3 + 2 = 6 。
# 第 2 本书不必放在第一层书架上。
