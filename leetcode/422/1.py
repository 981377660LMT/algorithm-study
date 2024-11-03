import time
from heapq import heappop, heappush
from random import randint
from typing import List, Tuple


MOD = int(1e9 + 7)
INF = int(1e20)


def max2(a: int, b: int) -> int:
    return a if a > b else b


DIR4 = [(0, 1), (0, -1), (1, 0), (-1, 0)]


class Solution:
    def minTimeToReach(self, moveTime: List[List[int]]) -> Tuple[int, int]:
        n, m = len(moveTime), len(moveTime[0])
        dist = [[[INF for _ in range(2)] for _ in range(m)] for _ in range(n)]
        dist[0][0][0] = 0
        queue = [(0, 0, 0, 0)]  # time, x, y, parity
        head = 0
        count = 0
        while head < len(queue):
            curTime, x, y, parity = queue[head]
            head += 1
            count += 1
            for dx, dy in DIR4:
                nx, ny = x + dx, y + dy
                if 0 <= nx < n and 0 <= ny < m:
                    cost = 1 if parity == 0 else 2
                    newTime = max2(curTime, moveTime[nx][ny]) + cost
                    nextParity = 1 ^ parity
                    if newTime < dist[nx][ny][nextParity]:
                        dist[nx][ny][nextParity] = newTime
                        queue.append((newTime, nx, ny, nextParity))
        return min(dist[n - 1][m - 1]), count

    def minTimeToReach2(self, moveTime: List[List[int]]) -> Tuple[int, int]:
        n, m = len(moveTime), len(moveTime[0])
        dist = [[[INF for _ in range(2)] for _ in range(m)] for _ in range(n)]
        dist[0][0][0] = 0
        pq = [(0, 0, 0, 0)]  # time, x, y, parity
        count = 0
        while pq:
            curTime, x, y, parity = heappop(pq)
            count += 1
            if x == n - 1 and y == m - 1:
                return curTime, count
            if curTime > dist[x][y][parity]:
                continue
            for dx, dy in DIR4:
                nx, ny = x + dx, y + dy
                if 0 <= nx < n and 0 <= ny < m:
                    cost = 1 if parity == 0 else 2
                    newTime = max2(curTime, moveTime[nx][ny]) + cost
                    nextParity = 1 - parity
                    if newTime < dist[nx][ny][nextParity]:
                        dist[nx][ny][nextParity] = newTime
                        heappush(pq, (newTime, nx, ny, nextParity))
        return -1, count


if __name__ == "__main__":
    s = Solution()

    n = 500
    m = 500
    moveTime = [[randint(0, int(1e9)) for _ in range(m)] for _ in range(n)]

    time1 = time.time()
    res1 = s.minTimeToReach(moveTime)
    time2 = time.time()
    print(time2 - time1, res1)

    time1 = time.time()
    res2 = s.minTimeToReach2(moveTime)
    time2 = time.time()
    print(time2 - time1, res2)
