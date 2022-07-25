from typing import List, Tuple, Optional
from collections import defaultdict, Counter


MOD = int(1e9 + 7)
INF = int(1e20)

# "Flush"：同花，五张相同花色的扑克牌。
# "Three of a Kind"：三条，有 3 张大小相同的扑克牌。
# "Pair"：对子，两张大小一样的扑克牌。
# "High Card"：高牌，五张大小互不相同的扑克牌。


class Solution:
    def bestHand(self, ranks: List[int], suits: List[str]) -> str:
        counter = Counter(ranks)
        if len(set(suits)) == 1:
            return "Flush"
        for v in counter.values():
            if v >= 3:
                return "Three of a Kind"
        for v in counter.values():
            if v >= 2:
                return "Pair"
        if len(set(ranks)) == 5:
            return "High Card"
        return ""
