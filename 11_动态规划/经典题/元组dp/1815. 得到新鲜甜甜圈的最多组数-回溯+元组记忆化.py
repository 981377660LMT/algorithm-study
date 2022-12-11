from functools import lru_cache
from typing import List, Tuple

# 有一个甜甜圈商店，每批次都烤 batchSize 个甜甜圈。
# 这个店铺有个规则，就是在烤一批新的甜甜圈时，之前 所有 甜甜圈都必须已经全部销售完毕。
# 给你一个整数 batchSize 和一个整数数组 groups ，
# 数组中的每个整数都代表一批前来购买甜甜圈的顾客，
# 其中 groups[i] 表示这一批顾客的人数。每一位顾客都恰好只要一个甜甜圈。

# 当有一批顾客来到商店时，他们所有人都必须在下一批顾客来之前购买完甜甜圈。
# 如果一批顾客中第一位顾客得到的甜甜圈不是上一组剩下的，那么这一组人都会很开心。
# 你可以随意安排每批顾客到来的顺序。请你返回在此前提下，最多 有多少组人会感到开心。

# 1 <= batchSize <= 9
# 1 <= groups.length <= 30
# 1 <= groups[i] <= 1e9

# 有点像1655. 分配重复整数


class Solution:
    def maxHappyGroups(self, batchSize: int, groups: List[int]) -> int:
        @lru_cache(None)
        def dfs(remain: int, mods: Tuple[int, ...]) -> int:
            """上一组剩下的个数 各个类型的组"""
            res, counter = 0, list(mods)
            for cur in range(1, batchSize):
                if mods[cur] > 0:  # 没有这个类型的组了
                    counter[cur] -= 1
                    res = max(res, (remain == 0) + dfs((remain - cur) % batchSize, tuple(counter)))
                    counter[cur] += 1
            return res

        modGroup = [0] * batchSize
        for g in groups:
            modGroup[g % batchSize] += 1
        res = dfs(0, tuple(modGroup)) + modGroup[0]
        dfs.cache_clear()
        return res
