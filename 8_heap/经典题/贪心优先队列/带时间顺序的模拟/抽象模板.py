# 模拟事件循环.性能待优化.


from heapq import heappop, heappush
import itertools
from typing import Any, Dict, Generator, List, Set, Tuple, TypeVar, Generic


R = TypeVar("R")


# —— 一、事件类型 —— #
class Event(Generic[R]):
    """基类：协程通过 `yield Event` 来告诉 Scheduler 我要等什么。"""

    pass


class Sleep(Event[R]):
    """睡眠事件：想要在 current_time + duration 时被唤醒。"""

    __slots__ = ("duration",)

    def __init__(self, duration: int):
        self.duration = duration


class Acquire(Event[R]):
    """请求资源，同时携带任意可比较的 priority。"""

    __slots__ = ("resource", "priority")

    def __init__(self, resource: R, priority: Any) -> None:
        self.resource = resource
        self.priority = priority


class Release(Event[R]):
    """释放资源事件：释放 `resource`，并立即尝试唤醒自己后续逻辑。"""

    __slots__ = ("resource",)

    def __init__(self, resource: R):
        self.resource = resource


# —— 二、Task & Scheduler —— #
# Generator 类型：yield 出 Event，send() 传入任意（我们 send Task 本身），不返回值
TaskGen = Generator[Event[R], None, None]


class Task(Generic[R]):
    def __init__(self, gen: TaskGen[R]):
        self.gen = gen


class Scheduler(Generic[R]):
    def __init__(self):
        self.current_time: int = 0
        # 时间队列：(wake_time, seq, Task)
        self._time_pq: List[Tuple[int, int, Task[R]]] = []
        # 资源等待：resource -> list of (priority_tuple, seq, Task)
        self._waiting: Dict[R, List[Tuple[Any, int, Task[R]]]] = {}
        # 资源占用标记
        self._locked: Set[R] = set()
        # 用于打破完全相同 priority 的场景
        self._counter = itertools.count()

    def start_task(self, task: Task[R]) -> None:
        ev = next(task.gen)
        self._dispatch(task, ev)

    def _dispatch(self, task: Task[R], ev: Event[R]) -> None:
        if isinstance(ev, Sleep):
            seq = next(self._counter)
            heappush(self._time_pq, (self.current_time + ev.duration, seq, task))

        elif isinstance(ev, Acquire):
            heap = self._waiting.setdefault(ev.resource, [])
            seq = next(self._counter)
            heappush(heap, (ev.priority, seq, task))

        elif isinstance(ev, Release):
            # 释放资源
            self._locked.discard(ev.resource)
            # 立即把同一个协程再跑下一步
            gen = task.gen
            try:
                ev2 = next(gen)
            except StopIteration:
                return
            self._dispatch(task, ev2)

        else:
            raise RuntimeError(f"Unknown event: {ev}")

    def run(self) -> None:
        # 一直跑，直到 time_pq 和 waiting 都空
        while self._time_pq or any(self._waiting.values()):
            # —— 1) 时间到的事件优先 ——
            if self._time_pq and self._time_pq[0][0] <= self.current_time:
                wake, _, task = heappop(self._time_pq)
                self.current_time = wake
                gen = task.gen
                try:
                    ev = next(gen)
                except StopIteration:
                    continue
                self._dispatch(task, ev)
                continue

            # —— 2) 资源可用，唤醒等待队列中优先级最高的 ——
            for res, heap in self._waiting.items():
                if res in self._locked:
                    continue
                if heap:
                    _, _, task = heappop(heap)
                    gen = task.gen
                    try:
                        ev = next(gen)
                    except StopIteration:
                        # 这条获取请求已经完成，资源仍保持 unlocked
                        continue
                    else:
                        # 这条请求真正拿到了资源
                        self._locked.add(res)
                        self._dispatch(task, ev)

            # —— 3) 如果都没就绪，就快进到下一个 wake_time ——
            if self._time_pq:
                self.current_time = self._time_pq[0][0]


# —— 三、Bridge 模拟业务 —— #
class Solution:
    def findCrossingTime(self, n: int, k: int, times: List[List[int]]) -> int:
        """
        LeetCode 2532: 用 k 个 worker 把 n 箱货从左岸运到右岸，
        每个 worker i 有 times[i] = [t1, t2, t3, t4]：
          - 左→右 过桥 t1
          - 右卸货   t2
          - 右→左 过桥 t3
          - 左装货   t4
        桥一次只能一个 worker，通过优先级 (side, -cross_time, -i) 决定谁先过。
        """
        scheduler = Scheduler[str]()
        remaining = n
        last_time = 0

        def worker(i: int) -> TaskGen[str]:
            # 拿到自己的 Task 对象
            nonlocal last_time, remaining

            cross = times[i][0] + times[i][2]
            while True:
                # —— 左岸申请过桥 —— side=1
                yield Acquire("bridge", (1, -cross, -i))
                # !先申请资源，获得资源后再检查是否还需要工作
                # # 如果已经没货了，挂起后直接退出
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

        # 启动所有 k 个 worker
        for i in range(k):
            t = Task(worker(i))
            scheduler.start_task(t)

        # 运行调度器
        scheduler.run()
        return last_time
