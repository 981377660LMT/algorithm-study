# 3885. 设计事件管理器
# https://leetcode.cn/problems/design-event-manager/description/
#
# 给你一组初始事件列表，其中每个事件有一个唯一的 eventId 和一个 priority（优先级）。
# 实现 EventManager 类：
# EventManager(int[][] events) 使用给定事件初始化管理器，其中 events[i] = [eventIdi, priorityi]。
# void updatePriority(int eventId, int newPriority) 更新具有 id 为 eventId 的 活跃 事件的优先级为 newPriority。
# int pollHighest() 移除并返回具有 最高优先级 的 活跃事件 的 eventId。如果有多个活动事件具有相同的优先级，则返回 eventId 最小的事件。如果没有活跃事件，则返回 -1。
# 如果一个事件没有被 pollHighest() 移除，则称其为 活跃事件。

from heapq import heapify, heappop, heappush


class EventManager:
    __slots__ = ("_idToPriority", "_pq")

    def __init__(self, events: list[list[int]]):
        self._idToPriority = {}
        self._pq = []
        for e, p in events:
            self._idToPriority[e] = p
            self._pq.append((-p, e))
        heapify(self._pq)

    def updatePriority(self, eventId: int, newPriority: int) -> None:
        self._idToPriority[eventId] = newPriority
        heappush(self._pq, (-newPriority, eventId))

    def pollHighest(self) -> int:
        while self._pq:
            p, e = heappop(self._pq)
            if self._idToPriority.get(e, -1) == -p:
                del self._idToPriority[e]
                return e
        return -1
