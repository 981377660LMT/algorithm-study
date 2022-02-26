from functools import lru_cache
from typing import List

MOD = int(1e9 + 7)

# 1105. 填充书架.
class Solution:
    def buildWall(self, height: int, width: int, bricks: List[int]) -> int:
        """相邻两行分割点不能重合"""

        @lru_cache(None)
        def dfs(index: int, preSplit: int, curSplit: int, curWidth: int) -> int:
            if index == height:
                return 1
            if curWidth == width:
                return dfs(index + 1, curSplit, 0, 0) % MOD

            res = 0
            for choose in bricks:
                nextWidth = curWidth + choose
                if nextWidth > width:
                    break

                nextSplit = curSplit
                if nextWidth != width:
                    nextSplit |= 1 << nextWidth
                if nextSplit & preSplit:
                    continue

                res += dfs(index, preSplit, nextSplit, nextWidth)
                res %= MOD

            return res % MOD

        bricks = sorted(bricks)
        res = dfs(0, 0, 0, 0)
        dfs.cache_clear()
        return res


print(Solution().buildWall(height=2, width=3, bricks=[1, 2]))
print(Solution().buildWall(height=1, width=1, bricks=[5]))
