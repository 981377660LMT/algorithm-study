# E - Packing Under Range Regulations
# !有编号为1,2,...,1e9的1e9个箱子和编号为1,2,...n的n个球
# !每个球只能放在范围[lefti,righti]的箱子里
# !球不能放在同一个位置，问可否放置成功？
# n<=2e5 1<=lefti<=righti<=1e9

# !按照时间遍历(注意不能顺序遍历,要跳着遍历)
# !1. 在每一个时间点，我们首先将当前时间点开始的会议加入小根堆
# !2. 如果会议已经结束,那么不可参加完所有会议
# !3.从剩下的会议中选择一个结束时间最早的去参加。

from heapq import heappop, heappush
import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(1e18)


def judgeEvents(events: List[List[int]]) -> bool:
    events.sort(key=lambda x: x[0])
    events.append([INF, INF])

    ei, time, pq = 0, 1, []
    while time < int(1e10):
        # !1. 在每一个时间点，我们首先将当前时间点开始的会议加入小根堆
        while ei < len(events) and events[ei][0] == time:
            heappush(pq, events[ei][1])
            ei += 1

        # !2. 如果会议已经结束,那么不可参加完所有会议
        if pq and pq[0] < time:
            return False

        # !3.从剩下的会议中选择一个结束时间最早的去参加。
        if pq:
            heappop(pq)
            time += 1
        else:
            time = events[ei][0]  # 加速遍历

    return True


if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        n = int(input())
        events = [list(map(int, input().split())) for _ in range(n)]
        res = judgeEvents(events)
        print("Yes" if res else "No")
# 5
# 2
# 1000000000 1000000000
# 1000000000 1000000000
# 2
# 999999999 1000000000
# 1000000000 1000000000
# 3
# 1000000000 1000000000
# 1000000000 1000000000
# 1000000000 1000000000
# 3
# 999999999 1000000000
# 1000000000 1000000000
# 999999998 1000000000
# 3
# 1000000000 1000000000
# 999999999 1000000000
# 999999999 1000000000
