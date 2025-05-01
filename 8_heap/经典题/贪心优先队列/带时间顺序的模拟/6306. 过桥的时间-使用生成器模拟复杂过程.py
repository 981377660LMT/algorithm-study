# 2532. 过桥的时间-使用生成器模拟复杂过程
# https://leetcode.cn/problems/time-to-cross-a-bridge/description/
#
# !本质上来说要模拟的是一个异步过程，每个worker都可以独立工作，但是共享“桥”这个资源，过桥的时候要加锁。
# 这个在现代编程中可以很容易用协程来表达，这里通过生成器来模拟协程。
# 核心是一个定时器堆和一个锁堆，定时器按照时间从小到大的顺序将事件依次取出执行，锁堆则按照题目要求将所有请求锁的worker进行排序，选出最优先的。
# 主循环每次首先将定时器中已经触发的协程继续执行一小段，直到它等待锁或者等待时间；然后如果桥没有被锁定，则取出等待桥的协程然后向前推进；如果都没有，则向前推进时间到下一个定时器事件。
# !具体的工作流程可以在worker协程中实现，这样无论是如何复杂的工作流程，都不影响主循环的逻辑。
#
# 传统做法可以用真正的协程/线程、锁和事件循环；在这里作者用 Python 生成器（generator）来模拟协程，并自己写了一个最简单的“事件调度器”——核心数据结构是两个堆（优先队列）：
# timer 堆：管理“未来某个时刻要被唤醒的协程”。
# bridge_lock 堆：管理“正在等桥”的协程队列，根据题目优先级规则选出下一个能过桥的工人。
# 主循环不断从这两个堆里弹事件，恢复对应的生成器执行“下一步”，直到所有货物运完。

from typing import List
from heapq import heappop, heappush


class Solution:
    def findCrossingTime(self, n: int, k: int, time: List[List[int]]) -> int:
        timer = []  # 定时器，按照时间从小到大的顺序将事件依次取出执行
        bridge_lock = []  # 锁堆，按照题目要求将所有请求锁的worker进行排序，选出最优先的
        bridge_in_use = False
        current_time = 0

        res = 0

        def wait_for_bridge(i, is_left, iter_self):
            nonlocal bridge_in_use
            # 工人请求使用桥 (加入等待队列)
            heappush(bridge_lock, (is_left, -time[i][0] - time[i][2], -i, iter_self))
            # 暂停等待调度器重新 next()
            yield
            # 被调度器选中，标记桥为占用状态
            bridge_in_use = True

        def wait_for_time(i, duration, iter_self):
            heappush(timer, (current_time + duration, i, iter_self))
            yield

        def worker(i):
            nonlocal n, res, bridge_in_use
            # 第一次 next() 时，会跑到这里接收自己的迭代器对象
            iter_self = yield

            while n > 0:
                # —— 在左侧申请过桥 ——
                yield from wait_for_bridge(i, True, iter_self)
                if n > 0:
                    # Go to right
                    yield from wait_for_time(i, time[i][0], iter_self)
                    bridge_in_use = False
                else:
                    bridge_in_use = False
                    break

                # 装货（每运出一件，n 减 1）
                n -= 1
                # 回去前卸货耗时 t2
                yield from wait_for_time(i, time[i][1], iter_self)
                # —— 在右侧申请过桥 ——
                yield from wait_for_bridge(i, False, iter_self)
                # 过桥回左侧耗时 t3
                yield from wait_for_time(i, time[i][2], iter_self)
                bridge_in_use = False

                # 记录最新送达时间
                res = current_time

                # 回左侧准备下一件耗时 t4
                yield from wait_for_time(i, time[i][3], iter_self)

        # -------------------- 主循环：事件驱动调度器 --------------------
        # 1. 启动所有 k 个 worker 协程，并让它们接收到自己的迭代器对象
        for i in range(k):
            iter_self = worker(i)
            next(iter_self)  # 跑到 yield 停下
            iter_self.send(iter_self)

        # 2. 不断处理 timer 和 lock 队列，直到所有事件完成
        while timer or bridge_lock:
            # —— 优先处理定时器事件 ——
            if timer and timer[0][0] <= current_time:
                _, _, next_iter = heappop(timer)
                try:
                    next(next_iter)  # 唤醒对应协程
                except StopIteration:
                    pass
            # —— 其次，如果桥空闲，就挑一个等待桥的协程运行它 ——
            elif bridge_lock and not bridge_in_use:
                _, _, _, next_iter = heappop(bridge_lock)
                try:
                    next(next_iter)
                except StopIteration:
                    pass
            # —— 如果都没就绪，就把时钟推进到下一个定时事件 ——
            else:
                current_time = timer[0][0]

        return res
