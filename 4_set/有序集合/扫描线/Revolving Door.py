from collections import defaultdict
from heapq import heapify, heappop
from typing import List

# 解析门


class Solution2:
    def solve(self, requests: List[List[int]]) -> List[List[int]]:
        """
        requests[i] contains [time, direction] meaning at time time,
        a person arrived at the door and either wanted to go in (1) or go out (0).

        先出后进，每次只能一个人进出，每次进出需要花费1时间
        输出最后每个人的[进出时间，进出方向]
        """
        res = []
        pq = requests
        heapify(pq)
        time, door = 0, 1
        while pq:
            curTime, direction = heappop(pq)
            if curTime > time:
                time = curTime
                door = 1
            res.append([time, door ^ 1 if direction == 1 else door])
            time += 1
            door = direction
        return res


print(Solution2().solve(requests=[[1, 0], [2, 1], [5, 0], [5, 1], [2, 0]]))
# The door starts as in
# [
#     [1, 0],
#     [2, 0],
#     [3, 1],
#     [5, 1],
#     [6, 0]
# ]
# At time 1, there's only one person so they can go out. Door becomes out.
# At time 2, there's two people but the person going out has priority so they go out.
# At time 3, the person looking to go in can now go in.
# At time 5, there's two people but the person going in has priority so they go out.
# At time 6, the last person can go out.
class Solution:
    def solve(self, requests):
        events = defaultdict(lambda: [0, 0])
        for key, direc in requests:
            events[key][direc] += 1

        res = []
        curTime = 0
        door = 1
        for key in sorted(events):
            if key > curTime:
                curTime = key
                door = 1

            count1 = events[key][door]  # 同向
            count2 = events[key][door ^ 1]  # 异向

            for key in range(curTime, curTime + count1):
                res.append([key, door])
            for key in range(curTime + count1, curTime + count1 + count2):
                res.append([key, door ^ 1])

            curTime += count1 + count2
            if count2:
                door ^= 1

        return res
