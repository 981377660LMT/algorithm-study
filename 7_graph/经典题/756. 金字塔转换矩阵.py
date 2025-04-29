# 756. 金字塔转换矩阵
# https://leetcode.cn/problems/pyramid-transition-matrix/solutions/911240/yi-ge-ke-yi-zha-diao-mu-qian-ji-hu-suo-y-6rwk/
# 砖块可以重用
# bottom 的长度范围在 [2, 8]。
# allowed 的长度范围在[0, 200]。
#
# !当前状态永远可以用上一层已经堆上的字母 + 下一层仍然露出来的字母来表示

from typing import List
from itertools import product
from functools import lru_cache


class Solution:
    def pyramidTransition(self, bottom: str, allowed: List[str]) -> bool:
        trans = dict()
        for p in allowed:
            trans.setdefault(p[:2], []).append(p[2])

        @lru_cache(None)
        def dfs(down: str, up: str):
            if len(up) >= 2:  # 剪枝
                if not dfs(up, ""):
                    return False
            if len(down) == 2:
                if not up:
                    return down in trans
                else:
                    return any(dfs(up + t, "") for t in trans.get(down, []))
            else:
                return any(dfs(down[1:], up + t) for t in trans.get(down[:2], []))

        return dfs(bottom, "")


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
