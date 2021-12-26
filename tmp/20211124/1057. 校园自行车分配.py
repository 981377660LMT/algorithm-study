from typing import List
from itertools import product

# 我们需要为每位工人分配一辆自行车
# 所有工人和自行车的位置都不相同。
# 在所有可用的自行车和工人中，我们选取彼此之间曼哈顿距离最短的工人自行车对  (worker, bike) ，并将其中的自行车分配給工人。
# 如果有多个 (worker, bike) 对之间的曼哈顿距离相同，那么我们选择工人索引最小的那对
# 如果有多种不同的分配方法，则选择自行车索引最小的一对


class Solution:
    def assignBikes(self, workers: List[List[int]], bikes: List[List[int]]) -> List[int]:
        m, n = len(workers), len(bikes)
        cands = []
        for i, j in product(range(m), range(n)):
            dis = abs(workers[i][0] - bikes[j][0]) + abs(workers[i][1] - bikes[j][1])
            cands.append((dis, i, j))

        res = [0] * m
        hit = 0
        s1, s2 = set(), set()
        for dis, i, j in sorted(cands):
            if i not in s1 and j not in s2:
                res[i] = j

                s1.add(i)
                s2.add(j)
                hit += 1
                if hit == m:
                    return res

        return res


print(Solution().assignBikes([[0, 0], [2, 1]], [[1, 2], [3, 3]]))
