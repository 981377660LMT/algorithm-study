from typing import List
from heapq import heappop, heappush

# 你有 k 个服务器，编号为 0 到 k-1
# 你的任务是找到 最繁忙的服务器 。最繁忙定义为一个服务器处理的请求数是所有服务器里最多的。
# 1882. 使用服务器处理任务


# 空闲的 cpu:优先条件为(index)，index小的先处理
# 忙碌的 cpu:优先条件为(endTime)，早结束早空闲
# i + (id - i) % k 表示下一个循环的编号


class Solution:
    def busiestServers(self, k: int, arrival: List[int], load: List[int]) -> List[int]:
        busy = []
        free = list(range(k))
        workTime = [0] * k
        for i, start in enumerate(arrival):
            # 1.busy执行完入队free
            while busy and busy[0][0] <= start:
                _, id = heappop(busy)
                heappush(free, i + (id - i) % k)
            # 2.从free取出一个server执行当前任务
            if free:
                id = heappop(free) % k
                heappush(busy, (start + load[i], id))
                workTime[id] += 1

        max_ = max(workTime)  # 注意这里要先提取出来 不要在生成式里算
        return [i for i in range(k) if workTime[i] == max_]


print(Solution().busiestServers(k=3, arrival=[1, 2, 3, 4, 5], load=[5, 2, 3, 3, 3]))
# 输出：[1]
# 解释：
# 所有服务器一开始都是空闲的。
# 前 3 个请求分别由前 3 台服务器依次处理。
# 请求 3 进来的时候，服务器 0 被占据，所以它呗安排到下一台空闲的服务器，也就是服务器 1 。
# 请求 4 进来的时候，由于所有服务器都被占据，该请求被舍弃。
# 服务器 0 和 2 分别都处理了一个请求，服务器 1 处理了两个请求。所以服务器 1 是最忙的服务器。
