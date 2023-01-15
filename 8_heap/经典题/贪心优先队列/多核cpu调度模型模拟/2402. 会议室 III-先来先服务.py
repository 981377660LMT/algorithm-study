# CPU调度-多线程
# !贪心优先队列/cpu调度/1882. 使用服务器处理任务.py


# 会议将会按以下方式分配给会议室：

# 1. 每场会议都会在未占用且编号 `最小` 的会议室举办。
# 2. 如果没有可用的会议室，会议将会`延期`，直到存在空闲的会议室。
#    延期会议的持续时间和原会议持续时间 `相同` 。
# 3. 当会议室处于未占用状态时，将会优先提供给`原 开始 时间`更早的会议。


# 总结:
# !两个pq来回倒
# !free:维护空闲的cpu，存储 (cpu)
# !busy:维护运行任务的cpu，存储 (endTime,cpu)，早结束早空闲

from heapq import heappop, heappush
from typing import List


class Solution:
    def mostBooked(self, n: int, meetings: List[List[int]]) -> int:
        """
        返回举办最多次会议的房间 编号 。

        如果存在多个房间满足此条件，则返回编号 最小 的房间。
        """
        counter = [0] * n
        meetings.sort(key=lambda x: x[0])
        free, busy = list(range(n)), []
        for start, end in meetings:
            # !1.busy里的任务结束了 归还CPU
            while busy and busy[0][0] <= start:
                _, cpu = heappop(busy)
                heappush(free, cpu)

            # !2. 需要一个CPU来执行任务
            if free:  # !有空闲的CPU
                cpu = heappop(free)
            else:  # !没有空闲的CPU 需要使用最早结束的CPU
                nextEnd, cpu = heappop(busy)
                end = nextEnd + (end - start)  # !延期执行
            counter[cpu] += 1
            heappush(busy, (end, cpu))

        return counter.index(max(counter))


if __name__ == "__main__":
    print(Solution().mostBooked(n=3, meetings=[[1, 20], [2, 10], [3, 5], [4, 9], [6, 8]]))  # 1
