# 2 <= n <= 28
# 1 <= firstPlayer < secondPlayer <= n

# 暴力即可
from functools import lru_cache
from itertools import product
from typing import List, Tuple

# 1900. 最佳运动员的比拼回合-元组为记忆化参数的dfs
# 1,2,3,4,5名选手 15配对 24配对 3晋级

# 后序dfs
# product每组取一个的巧妙使用
# Cartesian product of input iterables. Equivalent to nested for-loops.


class Solution:
    def earliestAndLatest(self, n: int, firstPlayer: int, secondPlayer: int) -> List[int]:
        @lru_cache(None)
        def dfs(players: Tuple[int, ...]) -> List[int]:
            n = len(players)
            pairs = [(players[i], players[~i]) for i in range((n + 1) // 2)]
            if (firstPlayer, secondPlayer) in pairs:
                return [1, 1]

            res = [int(1e20), -int(1e20)]

            # 每组选1个，笛卡尔积
            for winners in product(*pairs):
                if firstPlayer not in winners or secondPlayer not in winners:
                    continue
                nexts = dfs(tuple(sorted(winners)))
                res = [min(res[0], 1 + nexts[0]), max(res[1], 1 + nexts[1])]
            return res

        return dfs(tuple(range(1, n + 1)))


print(Solution().earliestAndLatest(n=11, firstPlayer=2, secondPlayer=4))
# 输出：[3,4]
