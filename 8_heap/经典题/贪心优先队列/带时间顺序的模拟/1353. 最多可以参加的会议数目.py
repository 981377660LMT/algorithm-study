from typing import List
from heapq import heappop, heappush

INF = int(4e18)

# 每个会议你至少参加一天
# 一天只能参加一个会议。

# !思路：排序；贪心，始终参加结束时间最早的会议

# 1 <= events.length <= 1e5
# events[i].length == 2
# 1 <= startDayi <= endDayi <= 1e5

# !模拟一个过程，而这个过程一般是按照 `时间顺序` 去执行:
# !0.分析状态与每种状态的优先级,用pq存储
# !1.弄清模拟的结束条件 (while ...)
# !2.每次while循环处理中的事件 : while一次性加入所有,if来处理不同event
#    !注意如果没有event要处理,则需要更新时间到下一次状态变化的时间


class Solution:
    def maxEvents(self, events: List[List[int]]) -> int:
        events.sort(key=lambda x: x[0])
        ei, res, pq = 0, 0, []
        curTime = 0
        while pq or ei < len(events):  # !所有的会议处理完毕，结束循环
            while ei < len(events) and events[ei][0] <= curTime:  # 1.将当前时间点开始的会议加入小根堆
                heappush(pq, events[ei][1])
                ei += 1
            while pq and pq[0] < curTime:  # 2.再把当前已经结束的会议移除出小根堆（因为已经无法参加了），
                heappop(pq)
            if pq:  # 3.从剩下的会议中选择一个结束时间最早的去参加。
                heappop(pq)
                res += 1
                curTime += 1
            else:
                if ei < len(events):
                    curTime = events[ei][0]  # !加速遍历
        return res


print(Solution().maxEvents([[1, 4], [4, 4], [2, 2], [3, 4], [1, 1]]))
