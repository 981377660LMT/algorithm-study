# low<=x<=high 范围内相交的直线数


# 直线在范围内可以描述成[左边的y值，右边的y值]
from collections import Counter
from sortedcontainers import SortedList


class Solution:
    def solve(self, lines, lo, hi):
        res = 0
        intervals = sorted((k * lo + b, k * hi + b) for k, b in lines)
        leftCounter = Counter(k * lo + b for k, b in lines)
        rightAbove = SortedList(k * hi + b for k, b in lines)
        rightBelow = SortedList()

        for left, right in intervals:
            rightAbove.remove(right)
            res += int(
                leftCounter[left] > 1
                or bool(rightAbove and rightAbove[0] <= right)
                or bool(rightBelow and rightBelow[-1] >= right)
            )
            rightBelow.add(right)

        return res


# 直线表示：y=kx+b

print(Solution().solve(lines=[[2, 3], [-3, 5], [4, 6]], lo=0, hi=1))

# 相交三种情形：
# 1.端点相同
# 2.左小右大
# 3.左大右小
