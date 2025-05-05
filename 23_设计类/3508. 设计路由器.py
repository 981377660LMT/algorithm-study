# 请你设计一个数据结构来高效管理网络路由器中的数据包。每个数据包包含以下属性：

# source：生成该数据包的机器的唯一标识符。
# destination：目标机器的唯一标识符。
# timestamp：该数据包到达路由器的时间戳。
# 实现 Router 类：

# Router(int memoryLimit)：初始化路由器对象，并设置固定的内存限制。

# memoryLimit 是路由器在任意时间点可以存储的 最大 数据包数量。
# 如果添加一个新数据包会超过这个限制，则必须移除 最旧的 数据包以腾出空间。
# bool addPacket(int source, int destination, int timestamp)：将具有给定属性的数据包添加到路由器。

# 如果路由器中已经存在一个具有相同 source、destination 和 timestamp 的数据包，则视为重复数据包。
# 如果数据包成功添加（即不是重复数据包），返回 true；否则返回 false。
# int[] forwardPacket()：以 FIFO（先进先出）顺序转发下一个数据包。

# 从存储中移除该数据包。
# 以数组 [source, destination, timestamp] 的形式返回该数据包。
# 如果没有数据包可以转发，则返回空数组。
# int getCount(int destination, int startTime, int endTime)：

# 返回当前存储在路由器中（即尚未转发）的，且目标地址为指定 destination 且时间戳在范围 [startTime, endTime]（包括两端）内的数据包数量。
# 注意：对于 addPacket 的查询会按照 timestamp 的递增顺序进行。


from bisect import bisect_left, bisect_right
from collections import defaultdict, deque
from typing import List, Tuple


class Router:
    __slots__ = ("_memoryLimit", "_packetQueue", "_packetSet", "_destToTimestamps")

    def __init__(self, memoryLimit: int):
        """
        初始化路由器对象，并设置固定的内存限制。
        memoryLimit 是路由器在任意时间点可以存储的 最大 数据包数量。
        如果添加一个新数据包会超过这个限制，则必须移除 最旧的 数据包以腾出空间。
        """
        self._memoryLimit = memoryLimit
        self._packetQueue = deque()
        self._packetSet = set()
        self._destToTimestamps = defaultdict(deque)  # 最好不要用deque，随机访问下标是O(n/64)的.

    def addPacket(self, source: int, destination: int, timestamp: int) -> bool:
        """
        将具有给定属性的数据包添加到路由器。
        如果路由器中已经存在一个具有相同 source、destination 和 timestamp 的数据包，则视为重复数据包。
        如果数据包成功添加（即不是重复数据包），返回 true；否则返回 false。

        !timestamp 按递增顺序给出.
        """

        packet = (source, destination, timestamp)
        if packet in self._packetSet:
            return False
        if len(self._packetQueue) == self._memoryLimit:
            self._popLeft()
        self._append(packet)
        return True

    def forwardPacket(self) -> List[int]:
        """
        以 FIFO（先进先出）顺序转发下一个数据包。
        从存储中移除该数据包。
        以数组 [source, destination, timestamp] 的形式返回该数据包。
        如果没有数据包可以转发，则返回空数组。
        """
        return self._popLeft()

    def getCount(self, destination: int, startTime: int, endTime: int) -> int:
        """
        返回当前存储在路由器中（即尚未转发）的，且目标地址为指定 destination 且时间戳在范围 [startTime, endTime]（包括两端）内的数据包数量。
        """
        times = self._destToTimestamps[destination]
        return bisect_right(times, endTime) - bisect_left(times, startTime)

    def _append(self, packet: Tuple[int, int, int]) -> None:
        self._packetQueue.append(packet)
        self._packetSet.add(packet)
        self._destToTimestamps[packet[1]].append(packet[2])

    def _popLeft(self) -> List[int]:
        if not self._packetQueue:
            return []
        res = self._packetQueue.popleft()
        self._packetSet.remove(res)
        self._destToTimestamps[res[1]].popleft()
        return res


# Your Router object will be instantiated and called as such:
# obj = Router(memoryLimit)
# param_1 = obj.addPacket(source,destination,timestamp)
# param_2 = obj.forwardPacket()
# param_3 = obj.getCount(destination,startTime,endTime)
