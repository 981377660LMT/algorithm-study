from collections import Counter, defaultdict
from typing import List, Tuple


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def findWinners(self, matches: List[List[int]]) -> List[List[int]]:
        all = set()
        win = Counter()
        lose = Counter()
        for winner, loser in matches:
            all.add(winner)
            all.add(loser)
            win[winner] += 1
            lose[loser] += 1

        allWin = all - set(lose.keys())
        lose1 = set([l for l, count in lose.items() if count == 1])
        return [sorted(allWin), sorted(lose1)]

