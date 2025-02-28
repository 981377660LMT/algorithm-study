# 1845. 座位预约管理系统
# 2336. 无限集中的最小数字
# !维护取消预约的座位，时间复杂度与n无关.
# https://leetcode.cn/problems/seat-reservation-manager/solutions/2838121/liang-chong-fang-fa-wei-hu-ke-yu-yue-de-tmub8/


from heapq import heappop, heappush


class SeatManager:
    __slots__ = "_seats", "_available"

    def __init__(self, n: int):
        self._seats = 0
        self._available = []

    def reserve(self) -> int:
        """返回最小的可用座位编号."""
        if self._available:
            return heappop(self._available)
        self._seats += 1
        return self._seats

    def unreserve(self, seatNumber: int) -> None:
        """释放座位编号."""
        heappush(self._available, seatNumber)
