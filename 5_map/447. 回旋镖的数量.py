from typing import List
from collections import Counter
from math import dist


# 给定平面上 n 对 互不相同 的点 points
# 回旋镖 是由点 (i, j, k) 表示的元组 ，
# 其中 i 和 j 之间的距离和 i 和 k 之间的欧式距离相等（需要考虑元组的顺序）。
# 返回平面上所有回旋镖的数量。


class Solution:
    def numberOfBoomerangs(self, points: List[List[int]]) -> int:
        res = 0
        for p2 in points:
            counter = Counter()
            for p1 in points:
                if p1 == p2:
                    continue
                counter[dist(p1, p2)] += 1
            res += sum(v * (v - 1) for v in counter.values())
        return res


print(Solution().numberOfBoomerangs([[0, 0], [1, 0], [2, 0]]))
print(Solution().numberOfBoomerangs([[0, 0], [1, 0], [-1, 0], [0, 1], [0, -1]]))

