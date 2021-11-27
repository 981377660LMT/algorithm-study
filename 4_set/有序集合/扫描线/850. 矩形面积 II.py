from typing import List, Sequence
from sortedcontainers import SortedList

MOD = 10 ** 9 + 7
LEFT = 0
RIGHT = 1
# 找出平面中所有矩形叠加覆盖后的总面积。 由于答案可能太大，请返回它对 10 ^ 9 + 7 取模的结果。
# 1.从左往右，线性扫描 注意排序 让x重合的左边界排在右边界前
# 2.扫描，看左边界还是有右边界
# 3.计算底边宽*高，高度使用有序容器维护


class Solution:
    def rectangleArea(self, rectangles: List[List[int]]) -> int:
        rec = []
        for x1, y1, x2, y2 in rectangles:
            rec.append((x1, LEFT, y1, y2))
            rec.append((x2, RIGHT, y1, y2))
        rec.sort()

        def calHeight(height_record: Sequence) -> int:
            height_sum = 0
            bound = -1
            for (lower, upper) in height_record:
                bound = max(bound, lower)
                height_sum += max(0, upper - bound)
                bound = max(bound, upper)
            return height_sum

        res = 0
        activate = SortedList()
        pre_x = rec[0][0]
        for x, bound, y1, y2 in rec:
            res += calHeight(activate) * (x - pre_x)
            res %= MOD
            if bound == LEFT:
                activate.add((y1, y2))
            else:
                activate.discard((y1, y2))
            pre_x = x
        return res


print(Solution().rectangleArea([[0, 0, 2, 2], [1, 0, 2, 3], [1, 0, 3, 1]]))
