from bisect import bisect_left
from typing import List

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def latestTimeCatchTheBus(self, buses: List[int], passengers: List[int], capacity: int) -> int:
        """返回你可以搭乘公交车的最晚到达公交站时间。你 不能 跟别的乘客同时刻到达。

        二分最晚时间+重叠元素应该最后处理
        """

        def check(mid: int) -> bool:
            """mid时能否上车 遍历公交车模拟过程 排序+遍历加指针记录"""
            pos = bisect_left(passengers, mid)
            queue = passengers[:pos] + [mid] + passengers[pos:]
            qi = 0
            for bt in buses:
                count = 0
                while qi < len(queue) and queue[qi] <= bt and count + 1 <= capacity:  # 当前qi能否上车
                    if qi == pos:  # 第pos个乘客上车了
                        return True
                    qi += 1
                    count += 1
            return False

        buses.sort()
        passengers.sort()
        bad = set(passengers)
        left, right = 1, int(1e10)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1

        res = right
        while res in bad:
            res -= 1
        return res


print(Solution().latestTimeCatchTheBus([3], [2, 4], 2))
