from typing import List
from functools import lru_cache

# 摆放书的顺序与你整理好的顺序相同。
# 以这种方式布置书架，返回书架整体可能的最小高度。
# 1 <= books.length <= 1000

# 考虑到之前的书能与现在组成一层
class Solution:
    def minHeightShelves(self, books: List[List[int]], shelfWidth: int) -> int:
        n = len(books)

        @lru_cache(None)
        def dfs(i, h, w) -> int:
            # dp(i)返回装到第i本书的最小高度,h记录本层最大高度,w记录当前层剩余宽度
            if i == n:
                return h

            nextH = max(h, books[i][1])
            if books[i][0] <= w:
                # 可以不用建新书架，也可以next Row新建书架
                res = min(
                    dfs(i + 1, nextH, w - books[i][0]),
                    h + dfs(i + 1, books[i][1], shelfWidth - books[i][0]),
                )
            else:
                # next Row新建书架
                res = h + dfs(i + 1, books[i][1], shelfWidth - books[i][0])

            return res

        return dfs(0, 0, shelfWidth)


print(
    Solution().minHeightShelves(
        books=[[1, 1], [2, 3], [2, 3], [1, 1], [1, 1], [1, 1], [1, 2]], shelfWidth=4
    )
)
# 输出：6
# 解释：
# 3 层书架的高度和为 1 + 3 + 2 = 6 。
# 第 2 本书不必放在第一层书架上。
