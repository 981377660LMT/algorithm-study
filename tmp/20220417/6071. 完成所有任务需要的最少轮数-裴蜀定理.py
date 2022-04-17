from collections import Counter
from typing import List, Tuple


MOD = int(1e9 + 7)
INF = int(1e20)


# 裴蜀定理
class Solution:
    def minimumRounds(self, tasks: List[int]) -> int:
        res = 0
        counter = Counter(tasks)

        for count in counter.values():
            if count == 1:
                return -1

            # 特判
            if count % 3 == 1:
                count -= 2
                res += 1

            div3, mod3 = divmod(count, 3)
            res += div3
            div2, mod2 = divmod(mod3, 2)
            res += div2
            if mod2:
                return -1

        return res

