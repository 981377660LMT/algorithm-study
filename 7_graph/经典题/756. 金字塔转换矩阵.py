from typing import List
from itertools import product
from collections import defaultdict
from functools import lru_cache

# 砖块可以重用
# bottom 的长度范围在 [2, 8]。
# allowed 的长度范围在[0, 200]。
class Solution:
    def pyramidTransition(self, bottom: str, allowed: List[str]) -> bool:
        cand = defaultdict(lambda: defaultdict(set))
        for left, right, up in allowed:
            cand[left][right].add(up)

        @lru_cache(None)
        def dfs(char: str) -> bool:
            if len(char) <= 1:
                return True
            # 可取组合的笛卡尔积
            for nextLevel in product(*(cand[left][right] for left, right in zip(char, char[1:]))):
                if dfs(nextLevel):
                    return True
            return False

        return dfs(bottom)


print(Solution().pyramidTransition(bottom="BCD", allowed=["BCG", "CDE", "GEA", "FFF"]))
# 输出：true
# 解释：
# 可以堆砌成这样的金字塔:
#     A
#    / \
#   G   E
#  / \ / \
# B   C   D

# 因为符合 BCG、CDE 和 GEA 三种规则。
print(*product(*([1, 2, 3], [4, 5])))

