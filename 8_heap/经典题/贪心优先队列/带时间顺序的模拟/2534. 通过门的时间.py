# 2534. 通过门的时间
# n 个人，按从 0 到 n - 1 编号。现在有一扇门，每个人只能通过门进入或离开一次，耗时一秒。
# 给你一个 非递减顺序 排列的整数数组 arrival ，数组长度为 n ，其中 arrival[i] 是第 i 个人到达门前的时间。
# 另给你一个长度为 n 的数组 state ，其中 state[i] 是 0 则表示第 i 个人希望进入这扇门，是 1 则表示 TA 想要离开这扇门。
# 如果 同时 有两个或更多人想要使用这扇门，则必须遵循以下规则：

# 如果前一秒 没有 使用门，那么想要 离开 的人会先离开。
# 如果前一秒使用门 进入 ，那么想要 进入 的人会先进入。
# 如果前一秒使用门 离开 ，那么想要 离开 的人会先离开。
# 如果多个人都想朝同一方向走（都进入或都离开），编号最小的人会先通过门。
# 返回一个长度为 n 的数组 answer ，其中 answer[i] 是第 i 个人通过门的时刻（秒）。


# !0 和 1 的模拟可以用数组/异或来代替分类讨论

from heapq import heappop, heappush
from typing import List


class Solution:
    def timeTaken(self, arrival: List[int], state: List[int]) -> List[int]:
        n = len(arrival)
        people = [(time, kind, id) for id, (time, kind) in enumerate(zip(arrival, state))]
        people.sort(key=lambda x: x[0])

        pq = [[], []]  # pq[0] for enter, pq[1] for exit
        remain, curTime = n, 0
        ei = 0
        preState = 1  # 0: enter, 1: exit
        res = [0] * n
        while remain > 0:
            while ei < n and people[ei][0] <= curTime:
                _, kind, id = people[ei]
                heappush(pq[kind], id)
                ei += 1
            if pq[0] or pq[1]:
                kind = preState if pq[preState] else preState ^ 1
                id = heappop(pq[kind])
                res[id] = curTime
                remain -= 1
                curTime += 1
                preState = kind
            elif ei < n:
                curTime = people[ei][0]
                preState = 1

        return res
