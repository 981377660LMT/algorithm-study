from functools import lru_cache
from typing import List, Tuple

# 1 <= batchSize <= 9
# 1 <= groups.length <= 30
# 1 <= groups[i] <= 109

# 有点像1655. 分配重复整数

# 1815. 得到新鲜甜甜圈的最多组数
class Solution:
    def maxHappyGroups(self, batchSize: int, groups: List[int]) -> int:
        @lru_cache(None)
        def bt(cur: int, mods: Tuple[int, ...]) -> int:
            res, counter = 0, list(mods)
            for m in range(batchSize):
                if counter[m] == 0:
                    continue
                counter[m] -= 1
                res = max(res, int(cur == 0) + bt((cur - m) % batchSize, tuple(counter)))
                counter[m] += 1

            return res

        mods = [0] * batchSize
        for g in groups:
            mods[g % batchSize] += 1

        return bt(0, tuple(mods))

