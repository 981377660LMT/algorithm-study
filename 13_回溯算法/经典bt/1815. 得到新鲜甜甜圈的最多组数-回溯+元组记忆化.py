from functools import lru_cache
from typing import List, Tuple

# 1 <= batchSize <= 9
# 1 <= groups.length <= 30
# 1 <= groups[i] <= 109

# 有点像1655. 分配重复整数


class Solution:
    def maxHappyGroups(self, batchSize: int, groups: List[int]) -> int:
        @lru_cache(None)
        def dfs(remain: int, mods: Tuple[int, ...]) -> int:
            """上一组剩下的个数 各个类型的组"""
            res = 0
            for cur in range(batchSize):
                if mods[cur] == 0:  # 没有这个类型的组了
                    continue
                nextMods = mods[:cur] + (mods[cur] - 1,) + mods[cur + 1 :]
                res = max(res, int(remain == 0) + dfs((remain - cur) % batchSize, tuple(nextMods)))

            return res

        modGroup = [0] * batchSize
        for g in groups:
            modGroup[g % batchSize] += 1
        return dfs(0, tuple(modGroup))

