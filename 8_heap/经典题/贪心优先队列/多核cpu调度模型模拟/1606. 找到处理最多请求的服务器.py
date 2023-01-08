from typing import List
from heapq import heappop, heappush

# 你有 k 个服务器，编号为 0 到 k-1
# 第 i （序号从 0 开始）个请求到达。
# !如果所有服务器都已被占据，那么该请求被舍弃（完全不处理）。
# 如果第 (i % k) 个服务器空闲，那么对应服务器会处理该请求。
# 否则，将请求安排给下一个空闲的服务器
# （服务器构成一个环，必要的话可能从第 0 个服务器开始继续找下一个空闲的服务器）。
# 比方说，如果第 i 个服务器在忙，那么会查看第 (i+1) 个服务器，第 (i+2) 个服务器等等。

# !你的任务是找到 最繁忙的服务器 。最繁忙定义为一个服务器处理的请求数是所有服务器里最多的。
# 1882. 使用服务器处理任务


# !两个堆来回倒
# !free:维护空闲的cpu，优先条件为(index)，index小的先处理
# !busy:维护运行任务的cpu，存储 (endTime,cpu)，早结束早空闲

# ! i + (cpu - i) % k 表示下一个循环的编号
# !这里是寻找下一个id；
# 如果cpu在i的顺时针方向，那么可以push cpu；
# 如果在逆时针方向，就要继续向走(cpu-i)%k (因为环要*2)


class Solution:
    def busiestServers(self, k: int, arrival: List[int], load: List[int]) -> List[int]:
        free, busy = list(range(k)), []
        workTime = [0] * k
        for i, (start, duration) in enumerate(zip(arrival, load)):
            # !1.busy里的任务结束了 归还CPU
            while busy and busy[0][0] <= start:
                _, cpu = heappop(busy)
                heappush(free, i + (cpu - i) % k)
            # !2. 需要一个CPU来执行任务
            if free:
                cpu = heappop(free) % k
                heappush(busy, (start + duration, cpu))
                workTime[cpu] += 1

        max_ = max(workTime)
        return [i for i in range(k) if workTime[i] == max_]


print(Solution().busiestServers(k=3, arrival=[1, 2, 3, 4, 5], load=[5, 2, 3, 3, 3]))
# 输出：[1]
# 解释：
# 所有服务器一开始都是空闲的。
# 前 3 个请求分别由前 3 台服务器依次处理。
# 请求 3 进来的时候，服务器 0 被占据，所以它呗安排到下一台空闲的服务器，也就是服务器 1 。
# 请求 4 进来的时候，由于所有服务器都被占据，该请求被舍弃。
# 服务器 0 和 2 分别都处理了一个请求，服务器 1 处理了两个请求。所以服务器 1 是最忙的服务器。
