"""
三种操作：

1.向A数组末尾添加一个数x
2.输出A数组第一个元素,并删除
3.对A数组进行升序排序

!不可能真的去排序 而是用一个优先队列模拟排序
"""

from collections import deque
from heapq import heappop, heappush
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    q = int(input())
    queue = deque()
    pq = []
    for _ in range(q):
        kind, *args = map(int, input().split())
        if kind == 1:
            x = args[0]
            queue.append(x)
        elif kind == 2:
            if pq:
                print(heappop(pq))
            else:
                print(queue.popleft())
        else:
            for num in queue:
                heappush(pq, num)
            queue.clear()
