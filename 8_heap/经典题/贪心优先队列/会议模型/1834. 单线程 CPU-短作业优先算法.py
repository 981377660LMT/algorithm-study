# CPU调度-单线程
# 现有一个单线程 CPU ，同一时间只能执行 最多一项 任务
# 如果 CPU 空闲，且任务队列中没有需要执行的任务，则 CPU 保持空闲状态。
# 如果 CPU 空闲，但任务队列中有需要执行的任务，则 CPU 将会选择 `执行时间最短` 的任务开始执行。
# 如果多个任务具有同样的最短执行时间，则选择下标最小的任务开始执行。
# 一旦某项任务开始执行，CPU 在 执行完整个任务 前都不会停止。
# CPU 可以在完成一项任务后，立即开始执行一项新任务。

# !返回 CPU 处理任务的顺序。
# !(短作业优先算法:pq维护(执行时间,任务编号))
# tasks.length == n
# 1 <= n <= 105
# 1 <= enqueueTimei, processingTimei <= 109

from heapq import heappop, heappush
from typing import List


class Solution:
    def getOrder(self, tasks: List[List[int]]) -> List[int]:
        """现有一个单线程 CPU ，同一时间只能执行 最多一项 任务

        Args:
            tasks (List[List[int]]):
            tasks[i] = [enqueueTimei, processingTimei] 意味着
            第 i 项任务将会于 enqueueTimei 时进入任务队列，
            需要 processingTimei 的时长完成执行。

        Returns:
            List[int]: 返回 CPU 处理任务的顺序。
        """
        res, n = [], len(tasks)
        events = sorted([(start, duration, i) for i, (start, duration) in enumerate(tasks)])

        ei = 0
        pq = []
        time = 0  # !因为不是先来先服务，所以要维护一个离散的时间戳
        while len(res) < n:
            # 1.在每一个时间点，我们首先将当前时间点开始的所有任务加入小根堆，
            while ei < n and events[ei][0] <= time:
                duration, taskId = events[ei][1], events[ei][2]
                heappush(pq, (duration, taskId))
                ei += 1
            # 2.从任务中选择一个最短的的去参加。
            if pq:
                duration, taskId = heappop(pq)
                time += duration
                res.append(taskId)
            elif ei < n:  # 3.如果没有任务，那么就跳到下一个任务开始的时间点
                time = events[ei][0]
        return res
