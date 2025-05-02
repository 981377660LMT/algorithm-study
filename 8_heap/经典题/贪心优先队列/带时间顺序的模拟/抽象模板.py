# https://leetcode.cn/problems/time-to-cross-a-bridge/solutions/3667441/python-mo-ni-shi-jian-xun-huan-by-981377-1ar6/
# 模拟事件循环.性能待优化.

import itertools
from heapq import heappop, heappush
from typing import Any, Generator, Hashable, List, Union


# region Events
Event = Union["Sleep", "Acquire", "Release"]


class Sleep:
    """睡眠事件：想要在 current_time + duration 时被唤醒。"""

    __slots__ = ("duration",)

    def __init__(self, duration: int):
        self.duration = duration


class Acquire:
    """请求资源，同时携带任意可比较的 priority。"""

    __slots__ = ("resource", "priority")

    def __init__(self, resource: Hashable, priority: Any):
        self.resource = resource
        self.priority = priority


class Release:
    """释放资源事件：释放 `resource`，并立即尝试唤醒自己后续逻辑。"""

    __slots__ = ("resource",)

    def __init__(self, resource: Hashable):
        self.resource = resource


# endregion

# region Scheduler & Task


TaskGen = Generator[Event, None, None]


class Task:
    """协程任务：携带一个生成器，生成器每次 yield 一个事件。"""

    __slots__ = ("gen",)

    def __init__(self, gen: TaskGen):
        self.gen = gen


class Scheduler:
    """调度器：负责管理所有的 Task 和资源的分配。"""

    __slots__ = ("current_time", "_time_pq", "_waiting", "_waiting_count", "_locked", "_counter")

    def __init__(self):
        self.current_time = 0
        # 时间队列：(wake_time, seq, Task)
        self._time_pq = []
        # 资源等待：resource -> list of (priority, seq, Task)
        self._waiting = {}
        self._waiting_count = 0
        # 资源占用标记
        self._locked = set()
        # 用于打破完全相同 priority 的场景
        self._counter = itertools.count()

    def start_task(self, task: Task) -> None:
        event = next(task.gen)
        self._dispatch(task, event)

    def run(self) -> None:
        # 一直跑，直到 time_pq 和 waiting 都空
        while self._time_pq or self._waiting_count:
            # —— 1) 时间到的事件优先 ——
            if self._time_pq and self._time_pq[0][0] <= self.current_time:
                wake, _, task = heappop(self._time_pq)
                self.current_time = wake
                event = next(task.gen)
                self._dispatch(task, event)
                continue

            # —— 2) 资源可用，唤醒等待队列中优先级最高的 ——
            for res, pq in self._waiting.items():
                if pq and res not in self._locked:
                    _, _, task = heappop(pq)
                    self._waiting_count -= 1
                    event = next(task.gen)
                    self._locked.add(res)
                    self._dispatch(task, event)

            # —— 3) 如果都没就绪，就快进到下一个 wake_time ——
            if self._time_pq:
                self.current_time = self._time_pq[0][0]

    def _dispatch(self, task: Task, event: Event) -> None:
        if isinstance(event, Sleep):
            seq = next(self._counter)
            heappush(self._time_pq, (self.current_time + event.duration, seq, task))

        elif isinstance(event, Acquire):
            pq = self._waiting.setdefault(event.resource, [])
            seq = next(self._counter)
            heappush(pq, (event.priority, seq, task))
            self._waiting_count += 1

        else:
            self._locked.discard(event.resource)

            # 立即把同一个协程再跑下一步
            try:
                ev2 = next(task.gen)
            except StopIteration:
                return
            self._dispatch(task, ev2)


# endregion

# region utils

# endregion


if __name__ == "__main__":

    class Solution:
        def findCrossingTime(self, n: int, k: int, times: List[List[int]]) -> int:
            """
            LeetCode 2532: 用 k 个 worker 把 n 箱货从左岸运到右岸，
            每个 worker i 有 times[i] = [t1, t2, t3, t4]:
              - 左→右 过桥 t1
              - 右卸货   t2
              - 右→左 过桥 t3
              - 左装货   t4
            桥一次只能一个 worker，通过优先级 (side, -cross_time, -i) 决定谁先过。
            """
            scheduler = Scheduler()
            remaining = n
            last_time = 0

            def worker(i: int):
                nonlocal last_time, remaining

                cross = times[i][0] + times[i][2]
                while True:
                    # —— 左岸申请过桥 —— side=1
                    yield Acquire("bridge", (1, -cross, -i))

                    # !先申请资源，获得资源后再检查是否还需要工作
                    # 如果已经没货了，挂起后直接退出
                    if remaining <= 0:
                        yield Release("bridge")
                        return

                    # 过桥
                    yield Sleep(times[i][0])
                    yield Release("bridge")

                    # 卸货
                    remaining -= 1
                    yield Sleep(times[i][1])

                    # —— 右岸申请回程 —— side=0
                    yield Acquire("bridge", (0, -cross, -i))
                    yield Sleep(times[i][2])
                    yield Release("bridge")

                    # 记录最后一箱到达时间
                    last_time = scheduler.current_time

                    # 回左岸装下一箱
                    yield Sleep(times[i][3])

            for i in range(k):
                t = Task(worker(i))
                scheduler.start_task(t)

            scheduler.run()
            return last_time
