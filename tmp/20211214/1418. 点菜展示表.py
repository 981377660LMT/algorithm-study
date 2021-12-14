from typing import List
from sortedcontainers import SortedSet
from collections import defaultdict, Counter

# 食物名=>set
# 每桌的情况=>每桌都是一个Counter


class Solution:
    def displayTable(self, orders: List[List[str]]) -> List[List[str]]:
        foods = SortedSet()
        tables = defaultdict(Counter)
        for _, tableId, food in orders:
            foods.add(food)
            tables[tableId][food] += 1

        return [["Table"] + [food for food in foods]] + [
            [tableId] + [str(tables[tableId][food]) for food in foods]
            for tableId in sorted(tables.keys(), key=int)
        ]


print(
    Solution().displayTable(
        orders=[
            ["David", "3", "Ceviche"],
            ["Corina", "10", "Beef Burrito"],
            ["David", "3", "Fried Chicken"],
            ["Carla", "5", "Water"],
            ["Carla", "5", "Ceviche"],
            ["Rous", "3", "Ceviche"],
        ]
    )
)

# 输出
# [
#     ["Table", "Beef Burrito", "Ceviche", "Fried Chicken", "Water"],
#     ["3", "0", "2", "1", "0"],
#     ["5", "0", "1", "0", "1"],
#     ["10", "1", "0", "0", "0"],
# ]

