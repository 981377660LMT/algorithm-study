from typing import List, Sequence, Tuple
from sortedcontainers import SortedList

MOD = int(1e9 + 7)

# 找出平面中所有矩形叠加覆盖后的总面积。 由于答案可能太大，请返回它对 10 ^ 9 + 7 取模的结果。
# 1.从左往右，线性扫描 注意排序 让x重合的左边界排在右边界前，先加后删
# 2.扫描，看左边界还是有右边界
# 3.计算底边宽*高，高度使用有序容器维护


class Solution:
    def rectangleArea(self, rectangles: List[List[int]]) -> int:
        """计算平面中所有 rectangles 所覆盖的 总面积"""
        #  O(N^2LogN)
        def calHeight(intervals: Sequence[Tuple[int, int]]) -> int:
            """"求有序区间的总覆盖长度"""
            if not intervals:
                return 0

            preStart, preEnd = intervals[0]
            res = preEnd - preStart
            for i in range(1, len(intervals)):
                curStart, curEnd = intervals[i]
                if curEnd <= preEnd:
                    continue
                elif curEnd > preEnd and curStart <= preEnd:
                    res += curEnd - preEnd
                    preEnd = curEnd
                else:
                    res += curEnd - curStart
                    preStart, preEnd = curStart, curEnd

            return res

        events = []
        for x1, y1, x2, y2 in rectangles:
            events.append((x1, -y1, -y2, 0))
            events.append((x2, y1, y2, 1))
        events.sort()

        res = 0
        sl = SortedList()
        preX = events[0][0]
        for x, y1, y2, flag in events:
            if x != preX:
                res += calHeight(sl) * (x - preX)
                res %= MOD
            if flag == 0:
                sl.add((-y1, -y2))
            else:
                sl.discard((y1, y2))
            preX = x
        return res


print(Solution().rectangleArea([[0, 0, 2, 2], [1, 0, 2, 3], [1, 0, 3, 1]]))
