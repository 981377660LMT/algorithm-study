from typing import List
from heapq import heappop, heappush, heapify

# 给你两个 下标从 0 开始 的整数数组 servers 和 tasks ，
# 长度分别为 n​​​​​​ 和 m​​​​​​ 。servers[i] 是第 i​​​​​​​​​​ 台服务器的 权重 ，
# 而 tasks[j] 是处理第 j​​​​​​ 项任务 所需要的时间（单位：秒）。
# 第 j 项任务在第 j 秒可以开始处理
# 处理第 j 项任务时，`你需要为它分配一台 `权重最小` 的空闲服务器。
# 如果存在多台相同权重的空闲服务器，请选择 下标最小 的服务器。`
# 如果一台空闲服务器在第 t 秒分配到第 j 项任务，那么在 t + tasks[j] 时它将恢复空闲状态。
# 返回res[j] 是第 j 项任务分配的服务器的下标。

# 总结:
# 两个pq来回倒
# !free:维护空闲的cpu，存储 (weight,index)
# !busy:维护运行任务的cpu，存储 (endTime,weight,index)，早结束早空闲


class Solution:
    def assignTasks(self, servers: List[int], tasks: List[int]) -> List[int]:
        """
        Args:
            servers (List[int]): servers[i] 是第 i 台服务器的 权重
            tasks (List[int]): tasks[j] 是处理第 j 项任务 所需要的时间（单位：秒）。

        Returns:
            List[int]: 构建长度为 m 的答案数组 res ，其中 res[j] 是第 j 项任务分配的服务器的下标。
        """
        free, busy = [(weight, cpu) for cpu, weight in enumerate(servers)], []
        heapify(free)

        res = []
        for start, duration in enumerate(tasks):
            # !1.busy里的任务结束了 归还CPU
            while busy and busy[0][0] <= start:
                _, weight, cpu = heappop(busy)
                heappush(free, (weight, cpu))

            # !2. 需要一个CPU来执行任务
            if free:  # !有空闲的CPU
                weight, cpu = heappop(free)
                heappush(busy, (start + duration, weight, cpu))
            else:  # !没有空闲的CPU 需要使用最早结束的CPU
                nextEnd, weight, cpu = heappop(busy)
                heappush(busy, (nextEnd + duration, weight, cpu))  # !延期执行
            res.append(cpu)

        return res


print(Solution().assignTasks(servers=[3, 3, 2], tasks=[1, 2, 3, 2, 1, 2]))
# 输入：servers = [3,3,2], tasks = [1,2,3,2,1,2]
# 输出：[2,2,0,2,1,2]
# 解释：事件按时间顺序如下：
# - 0 秒时，第 0 项任务加入到任务队列，使用第 2 台服务器处理到 1 秒。
# - 1 秒时，第 2 台服务器空闲，第 1 项任务加入到任务队列，使用第 2 台服务器处理到 3 秒。
# - 2 秒时，第 2 项任务加入到任务队列，使用第 0 台服务器处理到 5 秒。
# - 3 秒时，第 2 台服务器空闲，第 3 项任务加入到任务队列，使用第 2 台服务器处理到 5 秒。
# - 4 秒时，第 4 项任务加入到任务队列，使用第 1 台服务器处理到 5 秒。
# - 5 秒时，所有服务器都空闲，第 5 项任务加入到任务队列，使用第 2 台服务器处理到 7 秒。
