# https://www.cnblogs.com/Forever-666/p/14623282.html
# 雪灾与外卖

from heapq import heappop, heappush
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


def solve(
    n: int,
    m: int,
    mousePos: List[int],
    holePos: List[int],
    holeWage: List[int],
    holeCount: List[int],
) -> int:
    """雪灾与外卖

    Args:
        n (int): 老鼠的数量
        m (int): 洞的数量
        mousePos (List[int]): 每个老鼠的位置
        holePos (List[int]): 每个洞的位置
        holeWage (List[int]): 进入每个洞的花费
        holeCount (List[int]): 每个位置处洞的数量

    Returns:
        int: 求每个老鼠进洞的最小花费,如果洞的数量不够,则返回-1
    """
    events: List["Event"] = []
    count = 0
    for i in range(n):
        events.append(Event(mousePos[i], 0, 1, 0))
    for i in range(m):
        events.append(Event(holePos[i], 1, holeCount[i], holeWage[i]))
        count += holeCount[i]
    if count < n:
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


assert solve(4, 4, [8, 8, 9, 10], [1, 3, 5, 10], [4, 0, 0, 2], [1, 1, 1, 1]) == 22
