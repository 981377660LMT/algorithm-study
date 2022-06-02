# 2 <= n <= 28
# 1 <= firstPlayer < secondPlayer <= n

# 暴力即可
from functools import lru_cache
from itertools import combinations, product
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
                nexts = dfs(tuple(sorted(winners)))  # !注意这里要排序
                res = [min(res[0], 1 + nexts[0]), max(res[1], 1 + nexts[1])]
            return res

        return dfs(tuple(range(1, n + 1)))


print(Solution().earliestAndLatest(n=11, firstPlayer=2, secondPlayer=4))
# 输出：[3,4]

# 匹配足球比赛，给定1.国家和足球队的dict（e.g. {'Italy': [c1, c2]‍‍‌‍‌‍‍‍‍‌‍‍‌‍‍‌‍‍‌‍, 'Spain': [c3], 'France': [c4, c5, c6])
# 2. 之前踢过比赛的队[c3, c6], [c4, c2], [c5, c1]，要求匹配下一轮的n场比赛（任意一种组合），
# 要求同一个国家的不能一起比，之前踢过比赛的队也不能一起比。
mapping = {'C': ['a', 'b'], 'B': ['c', 'd']}
ok = [['b', 'd']]
visited = set([(a, b) for a, b in ok] + [(b, a) for a, b in ok])


def gen():
    for country1, country2 in combinations(mapping.keys(), 2):
        for team1, team2 in product(mapping[country1], mapping[country2]):
            if (team1, team2) not in visited:
                yield (team1, team2)
                visited.add((team1, team2))
                visited.add((team2, team1))


print(*gen())  # [('a', 'c'), ('a', 'd'), ('b', 'c')]
