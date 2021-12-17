from typing import List
from sortedcontainers import SortedList

# items[i] = [pricei, beautyi] 分别表示每一个物品的 价格 和 美丽值 。
# 对于每个查询 queries[j] ，你想求出价格小于等于 queries[j] 的物品中，
# 最大的美丽值 是多少。如果不存在符合条件的物品，那么查询的结果为 0 。

# query/items 按价格从小到大排序
# 不断更新availbel结构


class Solution:
    def maximumBeauty(self, items: List[List[int]], queries: List) -> List[int]:
        items = sorted(items)
        queries = sorted([(price, index) for index, price in enumerate(queries)])
        cur = SortedList()

        m, n = len(items), len(queries)
        ii, qi = 0, 0
        res = [0] * n

        while qi < n:
            # 添加阶段
            if ii < m and qi < n and items[ii][0] <= queries[qi][0]:
                cur.add(items[ii][1])
                ii += 1
            # 查询阶段
            else:
                if cur:
                    _, index = queries[qi]
                    res[index] = cur[-1]
                qi += 1

        return res


print(
    Solution().maximumBeauty(
        items=[[1, 2], [3, 2], [2, 4], [5, 6], [3, 5]], queries=[1, 2, 3, 4, 5, 6]
    )
)
# 输出：[2,4,5,5,6,6]
# 解释：
# - queries[0]=1 ，[1,2] 是唯一价格 <= 1 的物品。所以这个查询的答案为 2 。
# - queries[1]=2 ，符合条件的物品有 [1,2] 和 [2,4] 。
#   它们中的最大美丽值为 4 。
# - queries[2]=3 和 queries[3]=4 ，符合条件的物品都为 [1,2] ，[3,2] ，[2,4] 和 [3,5] 。
#   它们中的最大美丽值为 5 。
# - queries[4]=5 和 queries[5]=6 ，所有物品都符合条件。
#   所以，答案为所有物品中的最大美丽值，为 6 。
