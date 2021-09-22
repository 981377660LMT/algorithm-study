# 编写一个 O(1) 空间复杂度的一趟扫描算法，判断你所经过的路径是否相交。
from typing import List


class Solution:
    def isSelfCrossing(self, x: List[int]) -> bool:
        n = len(x)
        if n < 4:
            return False
        for i in range(3, n):
            if x[i] >= x[i - 2] and x[i - 1] <= x[i - 3]:
                return True
            if x[i - 1] <= x[i - 3] and x[i - 2] <= x[i]:
                return True
            if i > 3 and x[i - 1] == x[i - 3] and x[i] + x[i - 4] == x[i - 2]:
                return True
            if (
                i > 4
                and x[i] + x[i - 4] >= x[i - 2]
                and x[i - 1] >= x[i - 3] - x[i - 5]
                and x[i - 1] <= x[i - 3]
                and x[i - 2] >= x[i - 4]
                and x[i - 3] >= x[i - 5]
            ):
                return True
        return False
