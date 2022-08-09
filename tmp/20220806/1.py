from typing import List, Tuple, Optional
from collections import defaultdict, Counter


MOD = int(1e9 + 7)
INF = int(1e20)

# weighti 是所有价值为 valuei 物品的 重量之和 。


class Solution:
    def mergeSimilarItems(
        self, items1: List[List[int]], items2: List[List[int]]
    ) -> List[List[int]]:
        return sorted((Counter(dict(items1)) + Counter(dict(items2))).items())
