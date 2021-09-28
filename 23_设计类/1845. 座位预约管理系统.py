import heapq


class SeatManager:

    # 管理从 1 到 n 编号的 n 个座位。所有座位初始都是可预约的。
    def __init__(self, n: int):
        self.seats = list(range(1, n + 1))
        heapq.heapify(self.seats)

    # 返回可以预约座位的 最小编号 ，此座位变为不可预约。
    def reserve(self) -> int:
        return heapq.heappop(self.seats)

    # 将给定编号 seatNumber 对应的座位变成可以预约。
    def unreserve(self, seatNumber: int) -> None:
        heapq.heappush(self.seats, seatNumber)


# Your SeatManager object will be instantiated and called as such:
# obj = SeatManager(n)
# param_1 = obj.reserve()
# obj.unreserve(seatNumber)
