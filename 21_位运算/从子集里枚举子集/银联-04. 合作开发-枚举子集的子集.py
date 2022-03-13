from math import comb
from typing import Counter, List


MOD = int(1e9 + 7)
# 1 <= skills[i].length <= 4
# 2 <= skills.length <= 10^5
# 1 <= skills[i][j] <= 1000


class Solution:
    def coopDevelop(self, skills: List[List[int]]) -> int:
        n = len(skills)
        counter = Counter()
        for skill in skills:
            state = 0
            for s in skill:
                state |= 1 << (s - 1)
            counter[state] += 1

        res = comb(n, 2)
        for state in counter:
            tmp = 0
            cur = counter[state]
            g1 = state
            g1 = state & (g1 - 1)
            while g1:
                tmp = cur * counter[g1]
                tmp %= MOD
                g1 = state & (g1 - 1)
                res -= tmp
                res %= MOD

        # 要求真子集
        res -= sum(comb(c, 2) for c in counter.values())
        res %= MOD
        return res


print(Solution().coopDevelop([[1, 2, 3], [3], [2, 4]]))
print(Solution().coopDevelop([[3], [6]]))
print(Solution().coopDevelop([[2], [3, 5, 7], [2, 3, 5, 6], [3, 4, 8], [2, 6], [3, 4, 8], [3]]))
