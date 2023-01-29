from heapq import heappop, heappush
from random import randint
from time import time
from typing import List

INF = int(1e18)


class Event:
    __slots__ = ("pos", "kind", "count", "wage")

    def __init__(self, pos: int, kind: int, count: int, wage: int):
        self.pos = pos
        """位置"""
        self.kind = kind
        """类型 0:鼠 1:洞"""
        self.count = count
        """每个位置上的数量"""
        self.wage = wage
        """进入洞的花费"""


class Pair:
    __slots__ = ("value", "count")

    def __init__(self, value: int, count: int):
        self.value = value
        self.count = count

    def __lt__(self, other: "Pair"):
        return self.value < other.value


class Solution:
    def minimumTotalDistance(self, robot: List[int], factory: List[List[int]]) -> int:
        n, m = len(robot), len(factory)
        events: List["Event"] = []
        factoryCount = 0
        for i in range(n):
            pos = robot[i]
            events.append(Event(pos, 0, 1, 0))
        for i in range(m):
            pos, count = factory[i]
            events.append(Event(pos, 1, count, 0))
            factoryCount += count
        if factoryCount < n:
            return -1
        events.sort(key=lambda x: x.pos)

        mousePq: List["Pair"] = []
        holePq: List["Pair"] = []
        heappush(holePq, Pair(INF, INF))  # !一开始让所有老鼠都和一个距离它 INF 的洞匹配
        res = 0
        for event in events:
            if event.kind == 0:  # !老鼠
                top = heappop(holePq)
                res += top.value + event.pos
                heappush(mousePq, Pair(-event.pos * 2 - top.value, 1))
                top.count -= 1
                if top.count:
                    heappush(holePq, top)
            else:  # !洞
                rollback = 0
                while mousePq and event.count and mousePq[0].value + event.pos + event.wage < 0:
                    top = heappop(mousePq)
                    now = min(top.count, event.count)
                    res += now * (top.value + event.pos + event.wage)
                    heappush(holePq, Pair(-event.pos * 2 - top.value, now))
                    rollback += now
                    event.count -= now
                    top.count -= now
                    if top.count:
                        heappush(mousePq, top)
                if rollback:
                    heappush(mousePq, Pair(-event.pos - event.wage, rollback))
                if event.count:
                    heappush(holePq, Pair(-event.pos + event.wage, event.count))

        return res


assert Solution().minimumTotalDistance([1, 2, 3], [[2, 1], [3, 1], [4, 1]]) == 3
# 1e5 ranint
robots = [randint(1, int(1e5)) for _ in range(int(1e5))]
factories = [[randint(1, int(1e5)), randint(1, int(1e5))] for _ in range(int(1e5))]
time1 = time()
print(Solution().minimumTotalDistance(robots, factories))
time2 = time()
print(time2 - time1)
